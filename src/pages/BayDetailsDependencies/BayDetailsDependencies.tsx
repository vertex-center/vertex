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
import { SiHomebrew, SiSnapcraft } from "@icons-pack/react-simple-icons";
import { SegmentedButtons } from "../../components/SegmentedButton";

const pmIcons = {
    brew: <SiHomebrew />,
    snap: <SiSnapcraft />,
    sources: <Symbol name="folder_zip" />,
};

type Props = {
    name: string;
    dependency: DependencyModel;
    onChange: () => void;
};

export function Dependency(props: Props) {
    const { name, dependency, onChange } = props;

    const [packageManager, setPackageManager] = useState();

    const onPackageManagerChange = (pm: any) => setPackageManager(pm);

    const [installing, setInstalling] = useState(false);

    const install = () => {
        setInstalling(true);
        installDependencies([{ name, package_manager: packageManager }])
            .then(() => {
                onChange();
            })
            .catch(console.error)
            .finally(() => setInstalling(false));
    };

    return (
        <Horizontal alignItems="center" gap={16}>
            <div>{name}</div>
            <Spacer />
            {!dependency?.installed && (
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
            {dependency?.installed && (
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

    const reload = useCallback(() => {
        setIsLoading(true);
        getInstanceDependencies(uuid)
            .then((deps) => setDependencies(deps))
            .finally(() => setIsLoading(false));
    }, [uuid]);

    useEffect(() => {
        reload();
    }, [reload]);

    return (
        <Fragment>
            <Title>Dependencies</Title>
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
                {Object.entries(dependencies).map(([name, dep]) => (
                    <Dependency
                        key={name}
                        name={name}
                        dependency={dep}
                        onChange={reload}
                    />
                ))}
            </Vertical>
        </Fragment>
    );
}
