import { Fragment, useCallback, useEffect, useState } from "react";
import { Text, Title } from "../../../../components/Text/Text";
import { Horizontal, Vertical } from "../../../../components/Layouts/Layouts";
import Spacer from "../../../../components/Spacer/Spacer";
import Symbol from "../../../../components/Symbol/Symbol";
import styles from "./InstanceDependencies.module.sass";
import {
    Dependencies,
    Dependency as DependencyModel,
} from "../../../../models/dependency";
import { useParams } from "react-router-dom";
import classNames from "classnames";
import Button from "../../../../components/Button/Button";
import { SegmentedButtons } from "../../../../components/SegmentedButton";
import Progress from "../../../../components/Progress";
import Popup from "../../../../components/Popup/Popup";
import Code from "../../../../components/Code/Code";
import { api } from "../../../../backend/backend";
import { APIError } from "../../../../components/Error/Error";

type Props = {
    name: string;
    dependency: DependencyModel;
    onChange: () => void;
    onNeedsSudo: (command: string, packageManager: string) => void;
};

export function Dependency(props: Props) {
    const { name, dependency, onChange, onNeedsSudo } = props;

    const [packageManager, setPackageManager] = useState();

    const onPackageManagerChange = (pm: any) => setPackageManager(pm);

    const [installing, setInstalling] = useState(false);

    const install = () => {
        setInstalling(true);
        api.packages
            .install([{ name, package_manager: packageManager }])
            .then((res) => {
                if (res.data?.command) {
                    onNeedsSudo(res.data?.command, packageManager);
                } else {
                    onChange();
                }
            })
            .catch(console.error)
            .finally(() => setInstalling(false));
    };

    return (
        <Horizontal alignItems="center" gap={16}>
            <div>{name}</div>
            <Spacer />
            {!installing && !dependency?.installed && (
                <Horizontal alignItems="center" gap={12}>
                    Install with
                    <SegmentedButtons
                        value={packageManager}
                        onChange={onPackageManagerChange}
                        items={Object.keys(dependency?.install ?? {}).map(
                            (pm) => ({ value: pm })
                        )}
                        disabled={installing}
                    />
                    <Button
                        rightSymbol="download"
                        onClick={install}
                        loading={installing}
                        disabled={installing}
                    >
                        Install
                    </Button>
                </Horizontal>
            )}
            {installing && <Progress infinite />}
            {!installing && dependency?.installed && (
                <Horizontal
                    className={classNames({
                        [styles.installed]: dependency.installed,
                        [styles.notInstalled]: !dependency.installed,
                    })}
                    alignItems="center"
                    gap={4}
                >
                    <Symbol name="check" />
                    Installed
                </Horizontal>
            )}
        </Horizontal>
    );
}

export default function InstanceDependencies() {
    const { uuid } = useParams();

    const [dependencies, setDependencies] = useState<Dependencies>({});
    const [isLoading, setIsLoading] = useState(false);
    const [command, setCommand] = useState<string>();
    const [packageManager, setPackageManager] = useState<string>();
    const [error, setError] = useState();

    const reload = useCallback(() => {
        setIsLoading(true);
        api.instance.dependencies
            .get(uuid)
            .then((res) => setDependencies(res.data))
            .catch(setError)
            .finally(() => setIsLoading(false));
    }, [uuid]);

    useEffect(() => {
        reload();
    }, [reload]);

    const onNeedsSudo = (command: string, packageManager: string) => {
        setCommand(command);
        setPackageManager(packageManager);
    };

    const dismissPopup = () => {
        setCommand(undefined);
        setPackageManager(undefined);
    };

    return (
        <Fragment>
            <Title className={styles.title}>
                Dependencies ({Object.keys(dependencies).length})
            </Title>
            <APIError error={error} />
            <Horizontal alignItems="center">
                <Button
                    rightSymbol="refresh"
                    loading={isLoading}
                    disabled={isLoading}
                    onClick={reload}
                >
                    Reload
                </Button>
            </Horizontal>
            <Vertical className={styles.dependencies} gap={12}>
                {Object.entries(dependencies).map(([name, dep]) => (
                    <Dependency
                        key={name}
                        name={name}
                        dependency={dep}
                        onChange={reload}
                        onNeedsSudo={onNeedsSudo}
                    />
                ))}
            </Vertical>
            <Popup show={command !== undefined} onDismiss={dismissPopup}>
                <Title>Install dependencies</Title>
                <Text>
                    Vertex cannot install this dependency automatically, because{" "}
                    <code>{packageManager}</code> needs admin permissions.
                    Install manually using the command below, or try another
                    package manager.
                </Text>
                <Code
                    code={command ?? "FAILED TO RETRIEVE THE COMMAND."}
                    language={"bash"}
                />
                <Horizontal>
                    <Spacer />
                    <Button onClick={dismissPopup} primary rightSymbol="check">
                        OK
                    </Button>
                </Horizontal>
            </Popup>
        </Fragment>
    );
}
