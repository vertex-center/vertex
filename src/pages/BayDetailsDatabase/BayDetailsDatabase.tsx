import { Vertical } from "../../components/Layouts/Layouts";

import styles from "./BayDetailsDatabase.module.sass";
import { Title } from "../../components/Text/Text";
import {
    KeyValueGroup,
    KeyValueInfo,
} from "../../components/KeyValueInfo/KeyValueInfo";
import useInstance from "../../hooks/useInstance";
import { useParams } from "react-router-dom";
import InstanceSelect from "../../components/Input/InstanceSelect";
import { useState } from "react";
import { Instance } from "../../models/instance";
import Progress from "../../components/Progress";

export default function BayDetailsDatabase() {
    const { uuid } = useParams();
    const { instance } = useInstance(uuid);

    const [database, setDatabase] = useState<Instance>();

    const db = database?.service?.features?.databases?.[0];

    return (
        <Vertical gap={20}>
            <Title className={styles.title}>Database</Title>
            <Vertical gap={10}>
                {!instance && <Progress infinite />}
                {instance && (
                    <InstanceSelect
                        onChange={setDatabase}
                        instance={database}
                        query={{
                            features: instance?.service?.databases?.[0]?.types,
                        }}
                    />
                )}
                {db && (
                    <KeyValueGroup>
                        <KeyValueInfo name="Type" type="code">
                            {db?.type}
                        </KeyValueInfo>
                        <KeyValueInfo name="Port" type="code">
                            {db?.port}
                        </KeyValueInfo>
                        <KeyValueInfo name="Username" type="code">
                            {db?.username}
                        </KeyValueInfo>
                        <KeyValueInfo name="Password" type="code">
                            ***
                        </KeyValueInfo>
                    </KeyValueGroup>
                )}
            </Vertical>
        </Vertical>
    );
}
