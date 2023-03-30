import { Fragment, useEffect, useState } from "react";
import { Title } from "../../components/Text/Text";
import { Horizontal, Vertical } from "../../components/Layouts/Layouts";
import Spacer from "../../components/Spacer/Spacer";
import Symbol from "../../components/Symbol/Symbol";

import styles from "./BayDetailsDependencies.module.sass";
import {
    Dependencies,
    Dependency as DependencyModel,
    getInstanceDependencies,
} from "../../backend/backend";
import { useParams } from "react-router-dom";
import classNames from "classnames";
import Button from "../../components/Button/Button";
import { SiHomebrew, SiSnapcraft } from "@icons-pack/react-simple-icons";
import { SegmentedButtons } from "../../components/SegmentedButton/SegmentedButton";

const pmIcons = {
    brew: <SiHomebrew />,
    snap: <SiSnapcraft />,
    sources: <Symbol name="folder_zip" />,
};

type Props = {
    name: string;
    dependency: DependencyModel;
};

export function Dependency(props: Props) {
    const { name, dependency } = props;

    const [packageManager, setPackageManager] = useState();

    const onPackageManagerChange = (pm: any) => setPackageManager(pm);

    const install = () => {};

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
                    />
                    <Button rightSymbol="download" onClick={install}>
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

    useEffect(() => {
        getInstanceDependencies(uuid).then((deps) => setDependencies(deps));
    }, [uuid]);

    return (
        <Fragment>
            <Title>Dependencies</Title>
            <Vertical gap={12}>
                {Object.entries(dependencies).map(([name, dep]) => (
                    <Dependency name={name} dependency={dep} />
                ))}
            </Vertical>
        </Fragment>
    );
}
