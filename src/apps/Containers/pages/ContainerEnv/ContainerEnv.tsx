import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import {
    Button,
    Input,
    MaterialIcon,
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeadCell,
    TableRow,
    Title,
} from "@vertex-center/components";
import { Horizontal } from "../../../../components/Layouts/Layouts";
import { useContainerEnv } from "../../hooks/useContainer";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import Content from "../../../../components/Content/Content";
import { API } from "../../backend/api";
import { EnvVariables } from "../../backend/models";
import styles from "./ContainerEnv.module.sass";

export default function ContainerEnv() {
    const { uuid } = useParams();
    const queryClient = useQueryClient();

    const { env: currentEnv, isLoadingEnv, errorEnv } = useContainerEnv(uuid);
    const [env, setEnv] = useState<EnvVariables>();

    useEffect(() => {
        if (!currentEnv) return;
        setEnv(JSON.parse(JSON.stringify(currentEnv)));
        setSaved(undefined);
    }, [currentEnv]);

    const [saved, setSaved] = useState<boolean>(true);

    const mutationSaveEnv = useMutation({
        mutationFn: async (env: EnvVariables) => {
            await API.saveEnv(uuid, env);
        },
        onSuccess: () => setSaved(true),
        onSettled: () => {
            queryClient.invalidateQueries({
                queryKey: ["containers", uuid],
            });
            queryClient.invalidateQueries({
                queryKey: ["container_env", uuid],
            });
        },
    });
    const { isLoading: isUploading } = mutationSaveEnv;

    const save = () => mutationSaveEnv.mutate(env ?? []);

    const onChange = (i: number, e: React.ChangeEvent<HTMLInputElement>) => {
        const newEnv = [...env];
        newEnv[i].value = e.target.value;
        setEnv(newEnv);
        setSaved(isSaved());
    };

    const isSaved = () => {
        for (let i = 0; i < env.length; i++) {
            if (env[i].value !== currentEnv[i].value) {
                return false;
            }
        }
        return true;
    };

    return (
        <Content>
            <Title variant="h2">Environment</Title>
            <Table>
                <TableHead>
                    <TableRow>
                        <TableHeadCell>Name</TableHeadCell>
                        <TableHeadCell>Value</TableHeadCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {env?.map((env, i) => (
                        <TableRow key={env.name}>
                            <TableCell>
                                <Input
                                    id={env.name}
                                    value={env.name}
                                    className={styles.input}
                                    disabled
                                />
                            </TableCell>
                            <TableCell>
                                <Input
                                    id={env.name}
                                    value={env.value}
                                    name={env.name}
                                    placeholder={env.default}
                                    onChange={(v) => onChange(i, v)}
                                    type={env.secret ? "password" : undefined}
                                    disabled={isUploading}
                                    className={styles.input}
                                    style={{
                                        color:
                                            env.value !== currentEnv[i].value &&
                                            "var(--blue)",
                                    }}
                                />
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
            <ProgressOverlay show={isLoadingEnv ?? isUploading} />
            <Horizontal justifyContent="flex-end">
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
