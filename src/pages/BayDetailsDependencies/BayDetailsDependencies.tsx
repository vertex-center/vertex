import { Fragment, useCallback, useEffect, useState } from "react";
import { Title } from "../../components/Text/Text";
import { Horizontal, Vertical } from "../../components/Layouts/Layouts";
import Spacer from "../../components/Spacer/Spacer";
import Symbol from "../../components/Symbol/Symbol";

import styles from "./BayDetailsDependencies.module.sass";
import {
    Dependencies,
    Dependency as DependencyModel,
    getInstanceDependencies,
    installDependencies,
} from "../../backend/backend";
import { useParams } from "react-router-dom";
import classNames from "classnames";
import Button from "../../components/Button/Button";
import { SegmentedButtons } from "../../components/SegmentedButton";
import Progress from "../../components/Progress";
import { Error } from "../../components/Error/Error";

type Props = {
    name: string;
    dependency: DependencyModel;
    onChange: () => void;
    onNeedsSudo: (command: string) => void;
};

export function Dependency(props: Props) {
    const { name, dependency, onChange, onNeedsSudo } = props;

    const [packageManager, setPackageManager] = useState();

    const onPackageManagerChange = (pm: any) => setPackageManager(pm);

    const [installing, setInstalling] = useState(false);

    const install = () => {
        setInstalling(true);
        installDependencies([{ name, package_manager: packageManager }])
            .then((data: any) => {
                if (data?.command) {
                    onNeedsSudo(data?.command);
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

export default function BayDetailsDependencies() {
    const { uuid } = useParams();

    const [dependencies, setDependencies] = useState<Dependencies>({});
    const [isLoading, setIsLoading] = useState(false);
    const [command, setCommand] = useState<string>();

    const reload = useCallback(() => {
        setIsLoading(true);
        getInstanceDependencies(uuid)
            .then((deps) => setDependencies(deps))
            .finally(() => setIsLoading(false));
    }, [uuid]);

    useEffect(() => {
        reload();
    }, [reload]);

    const onNeedsSudo = (command: string) => {
        setCommand(command);
    };

    return (
        <Fragment>
            <Title>Dependencies ({Object.keys(dependencies).length})</Title>
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
            <Vertical gap={12}>
                {/* TODO: Show the command properly in a custom popup. */}
                {command && (
                    <Error
                        error={`This package manager needs sudo permissions. Install manually: ${command}`}
                    />
                )}
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
        </Fragment>
    );
}
