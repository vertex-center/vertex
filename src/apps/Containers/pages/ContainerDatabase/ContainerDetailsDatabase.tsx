import { Horizontal, Vertical } from "../../../../components/Layouts/Layouts";

import styles from "./ContainerDetailsDatabase.module.sass";
import { Title } from "../../../../components/Text/Text";
import {
    KeyValueGroup,
    KeyValueInfo,
} from "../../../../components/KeyValueInfo/KeyValueInfo";
import useContainer from "../../hooks/useContainer";
import { useParams } from "react-router-dom";
import ContainerSelect from "../../../../components/Input/ContainerSelect";
import { ChangeEvent, Fragment, useEffect, useState } from "react";
import { Container } from "../../../../models/container";
import Progress from "../../../../components/Progress";
import { Button, MaterialIcon } from "@vertex-center/components";
import { api } from "../../../../backend/api/backend";
import { DatabaseEnvironment } from "../../../../models/service";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import Spacer from "../../../../components/Spacer/Spacer";
import Input from "../../../../components/Input/Input";

type DatabaseProps = {
    container?: Container;

    dbID?: string;
    dbDefinition?: DatabaseEnvironment;

    onChange?: (name: string, dbUUID: string) => void;
};

function Database(props: Readonly<DatabaseProps>) {
    const { container, dbID, dbDefinition, onChange } = props;

    const [database, setDatabase] = useState<Container>();
    const [error, setError] = useState();

    const env = database?.environment;

    useEffect(() => {
        const uuid = container?.databases?.[dbID];
        if (uuid === undefined) return;
        api.vxContainers
            .container(uuid)
            .get()
            .then((data) => {
                setDatabase(data);
            })
            .catch(setError);
    }, [container]);

    const onDatabaseChange = (container: Container) => {
        setDatabase(container);
        onChange?.(dbID, container?.uuid);
    };

    const port = env?.[database?.service?.features?.databases?.[0]?.port];
    const username =
        env?.[database?.service?.features?.databases?.[0]?.username];

    if (error) return <APIError error={error} />;

    return (
        <Vertical gap={20}>
            <Title className={styles.title}>{dbDefinition?.display_name}</Title>
            <Vertical gap={10}>
                {!container && <Progress infinite />}
                {container && (
                    <ContainerSelect
                        onChange={onDatabaseChange}
                        container={database}
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

export default function ContainerDetailsDatabase() {
    const queryClient = useQueryClient();
    const { uuid } = useParams();
    const { container, isLoading, error } = useContainer(uuid);

    const [saved, setSaved] = useState<boolean>(undefined);

    const [databases, setDatabases] = useState<{
        [name: string]: {
            container_id: string;
            db_name?: string;
        };
    }>();

    const mutationSaveDatabase = useMutation({
        mutationFn: async () => {
            await api.vxContainers.container(uuid).patch({ databases });
        },
        onSuccess: () => {
            setSaved(true);
        },
        onSettled: () => {
            queryClient.invalidateQueries({
                queryKey: ["containers", uuid],
            });
        },
    });
    const { isLoading: isUploading, error: uploadingError } =
        mutationSaveDatabase;

    const onChange = (name: string, dbUUID: string) => {
        setDatabases((prev) => ({
            ...prev,
            [name]: {
                container_id: dbUUID,
            },
        }));
        setSaved(false);
    };

    const onChangeDbName = (
        e: ChangeEvent<HTMLInputElement>,
        db_name: string
    ) => {
        setDatabases((prev) => ({
            ...prev,
            [db_name]: {
                ...prev?.[db_name],
                db_name: e.target.value,
            },
        }));
        setSaved(false);
    };

    return (
        <Vertical gap={20}>
            {container &&
                Object.entries(container?.service?.databases ?? {}).map(
                    ([dbID, db]) => (
                        <Fragment key={dbID}>
                            <Database
                                dbID={dbID}
                                dbDefinition={db}
                                container={container}
                                onChange={onChange}
                            />
                            {databases?.[dbID]?.container_id && (
                                <Input
                                    label="Database name"
                                    onChange={(e: any) =>
                                        onChangeDbName(e, dbID)
                                    }
                                />
                            )}
                        </Fragment>
                    )
                )}
            <ProgressOverlay show={isLoading ?? isUploading} />
            <APIError error={error ?? uploadingError} />
            <Horizontal>
                <Spacer />
                <Button
                    variant="colored"
                    onClick={async () => mutationSaveDatabase.mutate()}
                    rightIcon={<MaterialIcon icon="save" />}
                    disabled={isUploading || saved || saved === undefined}
                >
                    Save{" "}
                    {container?.install_method === "docker" &&
                        "+ Recreate container"}
                </Button>
            </Horizontal>
        </Vertical>
    );
}
