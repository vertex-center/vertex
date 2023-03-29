import { Fragment, useEffect, useState } from "react";
import { Title } from "../../components/Text/Text";
import { Horizontal, Vertical } from "../../components/Layouts/Layouts";
import Spacer from "../../components/Spacer/Spacer";
import Symbol from "../../components/Symbol/Symbol";

import styles from "./BayDetailsDependencies.module.sass";
import { getInstance, Instance } from "../../backend/backend";
import { useParams } from "react-router-dom";

type Props = {
    dependency: string;
};

export function Dependency(props: Props) {
    const { dependency } = props;
    return (
        <Horizontal alignItems="center">
            <div>{dependency}</div>
            <Spacer />
            <Horizontal
                className={styles.installed}
                alignItems="center"
                gap={4}
            >
                <Symbol name="check" />
                Installed
            </Horizontal>
        </Horizontal>
    );
}

export default function BayDetailsDependencies() {
    const { uuid } = useParams();

    const [instance, setInstance] = useState<Instance>();

    useEffect(() => {
        getInstance(uuid).then((i: Instance) => setInstance(i));
    }, [uuid]);

    return (
        <Fragment>
            <Title>Dependencies</Title>
            <Vertical gap={12}>
                {Object.keys(instance?.dependencies ?? []).map((dep) => (
                    <Dependency dependency={dep} />
                ))}
            </Vertical>
        </Fragment>
    );
}
