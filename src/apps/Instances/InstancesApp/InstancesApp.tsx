import styles from "./InstancesApp.module.sass";
import Instance, { Instances } from "../../../components/Instance/Instance";
import { BigTitle } from "../../../components/Text/Text";
import { api } from "../../../backend/api/backend";
import { APIError } from "../../../components/Error/APIError";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useServerEvent } from "../../../hooks/useEvent";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import Toolbar from "../../../components/Toolbar/Toolbar";
import Spacer from "../../../components/Spacer/Spacer";
import Button from "../../../components/Button/Button";

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

    const toolbar = (
        <Toolbar>
            <Spacer />
            <Button primary to="/app/vx-instances/add" rightIcon="add">
                Create instance
            </Button>
        </Toolbar>
    );

    return (
        <div className={styles.server}>
            <ProgressOverlay show={isLoading} />
            <div className={styles.title}>
                <BigTitle>All instances</BigTitle>
            </div>

            {error && (
                <div className={styles.instances}>
                    <APIError error={error} />
                </div>
            )}

            {!isLoading && !isError && (
                <div className={styles.instances}>
                    {toolbar}
                    <Instances>
                        <Instance
                            instance={{
                                value: {
                                    display_name: "Vertex",
                                    status,
                                },
                            }}
                        />
                        {Object.values(instances)?.map((inst) => (
                            <Instance
                                key={inst.uuid}
                                instance={{
                                    value: inst,
                                    to: `/app/vx-instances/${inst.uuid}/`,
                                    onPower: async () =>
                                        mutationPower.mutate(inst.uuid),
                                }}
                            />
                        ))}
                    </Instances>
                </div>
            )}
        </div>
    );
}
