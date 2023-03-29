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

type Props = {
    name: string;
    dependency: DependencyModel;
};

export function Dependency(props: Props) {
    const { name, dependency } = props;
    return (
        <Horizontal alignItems="center">
            <div>{name}</div>
            <Spacer />
            <Horizontal
                className={classNames({
                    [styles.installed]: dependency.installed,
                    [styles.notInstalled]: !dependency.installed,
                })}
                alignItems="center"
                gap={4}
            >
                <Symbol name={dependency.installed ? "check" : "error"} />
                {dependency.installed ? "Installed" : "Not installed"}
            </Horizontal>
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
