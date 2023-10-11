import styles from "./InstancesApp.module.sass";
import Bay from "../../../components/Bay/Bay";
import { BigTitle } from "../../../components/Text/Text";
import Button from "../../../components/Button/Button";
import { Horizontal } from "../../../components/Layouts/Layouts";
import { api } from "../../../backend/backend";
import { APIError } from "../../../components/Error/APIError";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useServerEvent } from "../../../hooks/useEvent";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export default function InstancesApp() {
    const queryClient = useQueryClient();

    const queryInstances = useQuery({
        queryKey: ["instances"],
        queryFn: api.vxInstances.instances.all,
    });
    const { data: instances, isLoading, isError, error } = queryInstances;

    const mutationPower = useMutation({
        mutationFn: async (uuid: string) => {
            if (
                instances[uuid].status === "off" ||
                instances[uuid].status === "error"
            ) {
                await api.vxInstances.instance(uuid).start();
            } else {
                await api.vxInstances.instance(uuid).stop();
            }
        },
    });

    let status = "Waiting server response...";
    if (queryInstances.isSuccess) {
        status = "running";
    } else if (queryInstances.isError) {
        status = "off";
    }

    useServerEvent("/app/vx-instances/instances/events", {
        change: () => {
            queryClient.invalidateQueries({
                queryKey: ["instances"],
            });
        },
        open: () => {
            queryClient.invalidateQueries({
                queryKey: ["instances"],
            });
        },
    });

    return (
        <div className={styles.server}>
            <ProgressOverlay show={isLoading} />
            <div className={styles.title}>
                <BigTitle>All instances</BigTitle>
            </div>

            {error && (
                <div className={styles.bays}>
                    <APIError error={error} />
                </div>
            )}

            {!isLoading && !isError && (
                <div className={styles.bays}>
                    <Bay
                        instances={[
                            {
                                value: {
                                    display_name: "Vertex",
                                    status,
                                },
                            },
                        ]}
                    />
                    {Object.values(instances)?.map((inst) => (
                        <Bay
                            key={inst.uuid}
                            instances={[
                                {
                                    value: inst,
                                    to: `/app/vx-instances/${inst.uuid}/`,
                                    onPower: async () =>
                                        mutationPower.mutate(inst.uuid),
                                },
                            ]}
                        />
                    ))}
                    <Horizontal className={styles.addBay} gap={10}>
                        <Button
                            primary
                            to="/app/vx-instances/add"
                            leftIcon="add"
                        >
                            Create instance
                        </Button>
                    </Horizontal>
                </div>
            )}
        </div>
    );
}