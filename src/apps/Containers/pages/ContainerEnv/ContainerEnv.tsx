import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import EnvVariableInput from "../../components/EnvVariableInput/EnvVariableInput";
import {
    Button,
    InlineCode,
    MaterialIcon,
    Table,
    Title,
} from "@vertex-center/components";
import { Horizontal } from "../../../../components/Layouts/Layouts";
import { useContainerEnv } from "../../hooks/useContainer";
import styles from "./ContainerEnv.module.sass";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import Content from "../../../../components/Content/Content";
import { API } from "../../backend/api";
import { EnvVariables } from "../../backend/models";

export default function ContainerEnv() {
    const { uuid } = useParams();
    const queryClient = useQueryClient();

    const { env: currentEnv, isLoadingEnv, errorEnv } = useContainerEnv(uuid);
    const [env, setEnv] = useState<EnvVariables>(currentEnv);

    useEffect(() => {
        setEnv(currentEnv);
        setSaved(undefined);
    }, [currentEnv]);

    // undefined = not saved AND never modified
    const [saved, setSaved] = useState<boolean>(undefined);

    const mutationSaveEnv = useMutation({
        mutationFn: async (env: EnvVariables) => {
            await API.saveEnv(uuid, env);
        },
        onSuccess: () => setSaved(true),
        onSettled: () => {
            queryClient.invalidateQueries({
                queryKey: ["containers", uuid],
            });
        },
    });
    const { isLoading: isUploading } = mutationSaveEnv;

    const save = () => mutationSaveEnv.mutate(env ?? []);

    const onChange = (i: number, value: any) => {
        const newEnv = [...env];
        newEnv[i].value = value;
        setEnv(newEnv);
        setSaved(false);
    };

    return (
        <Content>
            <Title variant="h2">Environment</Title>
            <Table>
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Value</th>
                    </tr>
                </thead>
                <tbody>
                    {env?.map((env, i) => (
                        <tr key={env.name}>
                            <td>
                                <InlineCode>{env.name}</InlineCode>
                            </td>
                            <td>
                                <EnvVariableInput
                                    id={env.name}
                                    env={env}
                                    value={env.value}
                                    onChange={(v) => onChange(i, v)}
                                    disabled={isUploading}
                                />
                            </td>
                        </tr>
                    ))}
                </tbody>
            </Table>
            <ProgressOverlay show={isLoadingEnv ?? isUploading} />
            <Horizontal justifyContent="flex-end">
                {saved && (
                    <Horizontal
                        className={styles.saved}
                        alignItems="center"
                        gap={4}
                    >
                        <MaterialIcon icon="check" />
                        Saved!
                    </Horizontal>
                )}
                <Button
                    variant="colored"
                    onClick={save}
                    rightIcon={<MaterialIcon icon="save" />}
                    disabled={isUploading || saved || saved === undefined}
                >
                    Save
                </Button>
            </Horizontal>
            <APIError error={errorEnv} />
        </Content>
    );
}
