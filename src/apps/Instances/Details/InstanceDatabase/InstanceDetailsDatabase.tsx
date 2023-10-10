import { Vertical } from "../../../../components/Layouts/Layouts";

import styles from "./InstanceDetailsDatabase.module.sass";
import { Title } from "../../../../components/Text/Text";
import {
    KeyValueGroup,
    KeyValueInfo,
} from "../../../../components/KeyValueInfo/KeyValueInfo";
import useInstance from "../../../../hooks/useInstance";
import { useParams } from "react-router-dom";
import InstanceSelect from "../../../../components/Input/InstanceSelect";
import { useEffect, useState } from "react";
import { Instance } from "../../../../models/instance";
import Progress from "../../../../components/Progress";
import Button from "../../../../components/Button/Button";
import { api } from "../../../../backend/backend";
import { DatabaseEnvironment } from "../../../../models/service";
import { APIError } from "../../../../components/Error/APIError";

type DatabaseProps = {
    instance?: Instance;

    dbID?: string;
    dbDefinition?: DatabaseEnvironment;

    onChange?: (name: string, dbUUID: string) => void;
};

function Database(props: Readonly<DatabaseProps>) {
    const { instance, dbID, dbDefinition, onChange } = props;

    const [database, setDatabase] = useState<Instance>();
    const [error, setError] = useState();

    const env = database?.environment;

    useEffect(() => {
        const uuid = instance?.databases?.[dbID];
        if (uuid === undefined) return;
        api.vxInstances
            .instance(uuid)
            .get()
            .then((res) => {
                setDatabase(res.data);
            })
            .catch(setError);
    }, [instance]);

    const onDatabaseChange = (instance: Instance) => {
        setDatabase(instance);
        onChange?.(dbID, instance?.uuid);
    };

    const port = env?.[database?.service?.features?.databases?.[0]?.port];
    const username =
        env?.[database?.service?.features?.databases?.[0]?.username];

    if (error) return <APIError error={error} />;

    return (
        <Vertical gap={20}>
            <Title className={styles.title}>{dbDefinition?.display_name}</Title>
            <Vertical gap={10}>
                {!instance && <Progress infinite />}
                {instance && (
                    <InstanceSelect
                        onChange={onDatabaseChange}
                        instance={database}
                        query={{
                            features: dbDefinition?.types,
                        }}
                    />
                )}
                {database && (
                    <KeyValueGroup>
                        <KeyValueInfo name="Type" type="code">
                            {database?.service?.features?.databases?.[0]?.type}
                        </KeyValueInfo>
                        <KeyValueInfo name="Port" type="code">
                            {port}
                        </KeyValueInfo>
                        <KeyValueInfo name="Username" type="code">
                            {username}
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

export default function InstanceDetailsDatabase() {
    const { uuid } = useParams();
    const { instance } = useInstance(uuid);

    const [saved, setSaved] = useState<boolean>(undefined);
    const [uploading, setUploading] = useState<boolean>(false);

    const [databases, setDatabases] = useState<{
        [name: string]: string;
    }>();

    const save = () => {
        setUploading(true);
        api.vxInstances
            .instance(uuid)
            .patch({ databases })
            .then(() => setSaved(true))
            .catch(console.error)
            .finally(() => setUploading(false));
    };

    const onChange = (name: string, dbUUID: string) => {
        console.log({ ...databases, [name]: dbUUID });
        setDatabases((prev) => ({ ...prev, [name]: dbUUID }));
        setSaved(false);
    };

    return (
        <Vertical gap={20}>
            {instance &&
                Object.entries(instance?.service?.databases ?? {}).map(
                    ([dbID, db]) => (
                        <Database
                            key={dbID}
                            dbID={dbID}
                            dbDefinition={db}
                            instance={instance}
                            onChange={onChange}
                        />
                    )
                )}
            <Button
                primary
                large
                onClick={save}
                rightIcon="save"
                loading={uploading}
                disabled={saved || saved === undefined}
            >
                Save{" "}
                {instance?.install_method === "docker" &&
                    "+ Recreate container"}
            </Button>
        </Vertical>
    );
}
