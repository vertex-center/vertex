import { Fragment, useEffect, useState } from "react";
import { Title } from "../../../../components/Text/Text";
import { useParams } from "react-router-dom";
import EnvVariableInput from "../../../../components/EnvVariableInput/EnvVariableInput";
import { Button, MaterialIcon } from "@vertex-center/components";
import { Horizontal } from "../../../../components/Layouts/Layouts";
import useContainer from "../../hooks/useContainer";
import { Env, EnvVariable } from "../../../../models/service";
import styles from "./ContainerEnv.module.sass";
import { api } from "../../../../backend/api/backend";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export default function ContainerEnv() {
    const { uuid } = useParams();
    const queryClient = useQueryClient();

    const [env, setEnv] = useState<
        {
            env: EnvVariable;
            value: any;
        }[]
    >();

    const { container, isLoading, error } = useContainer(uuid);

    // undefined = not saved AND never modified
    const [saved, setSaved] = useState<boolean>(undefined);

    const mutationSaveEnv = useMutation({
        mutationFn: async (env: Env) => {
            await api.vxContainers.container(uuid).env.save(env);
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
    const { isLoading: isUploading } = mutationSaveEnv;

    const save = () => {
        const _env: Env = {};
        env.forEach((e) => {
            _env[e.env.name] = e.value;
        });
        mutationSaveEnv.mutate(_env);
    };

    useEffect(() => {
        setEnv(
            container?.service?.environment?.map((e) => ({
                env: e,
                value: container?.environment[e.name] ?? e.default ?? "",
            }))
        );
    }, [container]);

    const onChange = (i: number, value: any) => {
        setEnv((prev) =>
            prev.map((el, index) => {
                if (index !== i) return el;
                return { ...el, value };
            })
        );
        setSaved(false);
    };

    return (
        <Fragment>
            <ProgressOverlay show={isLoading ?? isUploading} />
            <Title className={styles.title}>Environment</Title>
            {env?.map((env, i) => (
                <EnvVariableInput
                    key={env.env.name}
                    env={env.env}
                    value={env.value}
                    onChange={(v) => onChange(i, v)}
                    disabled={isUploading}
                />
            ))}
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
                    Save{" "}
                    {container?.install_method === "docker" &&
                        "+ Recreate container"}
                </Button>
            </Horizontal>
            <APIError error={error} />
        </Fragment>
    );
}
