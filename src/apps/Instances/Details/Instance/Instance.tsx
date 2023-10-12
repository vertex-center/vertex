import Bay from "../../../../components/Bay/Bay";
import { useState } from "react";
import { api } from "../../../../backend/api/backend";
import { Outlet, useNavigate, useParams } from "react-router-dom";

import styles from "./Instance.module.sass";
import { Horizontal } from "../../../../components/Layouts/Layouts";
import Spacer from "../../../../components/Spacer/Spacer";
import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../../../components/Sidebar/Sidebar";
import Popup from "../../../../components/Popup/Popup";
import { Text, Title } from "../../../../components/Text/Text";
import Button from "../../../../components/Button/Button";
import Progress from "../../../../components/Progress";
import { SiDocker } from "@icons-pack/react-simple-icons";
import useInstance from "../../../../hooks/useInstance";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useServerEvent } from "../../../../hooks/useEvent";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Instance as InstanceModel } from "../../../../models/instance";

export default function Instance() {
    const { uuid } = useParams();
    const navigate = useNavigate();
    const queryClient = useQueryClient();

    const { instance, isLoading } = useInstance(uuid);

    const [showDeletePopup, setShowDeletePopup] = useState<boolean>();

    const route = uuid ? `/app/vx-instances/instance/${uuid}/events` : "";

    useServerEvent(route, {
        status_change: (e) => {
            queryClient.setQueryData(
                ["instances", uuid],
                (instance: InstanceModel) => ({ ...instance, status: e.data })
            );
        },
    });

    const mutationInstancePower = useMutation({
        mutationFn: async () => {
            if (instance.status === "off" || instance.status === "error") {
                await api.vxInstances.instance(uuid).start();
            } else {
                await api.vxInstances.instance(uuid).stop();
            }
        },
    });

    const mutationDeleteInstance = useMutation({
        mutationFn: api.vxInstances.instance(uuid).delete,
        onSuccess: () => {
            navigate("/app/vx-instances");
        },
    });
    const { isLoading: isDeleting, error: errorDeleting } =
        mutationDeleteInstance;

    const dismissDeletePopup = () => {
        setShowDeletePopup(false);
    };

    const content = (
        <Horizontal className={styles.content}>
            <Sidebar root={`/app/vx-instances/${uuid}`}>
                <SidebarGroup>
                    <SidebarItem
                        to={`/app/vx-instances/${uuid}/home`}
                        icon="home"
                        name="Home"
                    />
                </SidebarGroup>
                <SidebarGroup title="Analyze">
                    <SidebarItem
                        to={`/app/vx-instances/${uuid}/logs`}
                        icon="terminal"
                        name="Logs"
                    />
                    {instance?.install_method === "docker" && (
                        <SidebarItem
                            to={`/app/vx-instances/${uuid}/docker`}
                            icon={<SiDocker size={20} />}
                            name="Docker"
                        />
                    )}
                </SidebarGroup>
                <SidebarGroup title="Manage">
                    <SidebarItem
                        to={`/app/vx-instances/${uuid}/environment`}
                        icon="tune"
                        name="Environment"
                    />
                    {instance?.service?.databases && (
                        <SidebarItem
                            to={`/app/vx-instances/${uuid}/database`}
                            icon="database"
                            name="Database"
                        />
                    )}
                    <SidebarItem
                        to={`/app/vx-instances/${uuid}/update`}
                        icon="update"
                        name="Update"
                        notifications={
                            instance?.service_update?.available ? 1 : undefined
                        }
                    />
                    <SidebarItem
                        to={`/app/vx-instances/${uuid}/settings`}
                        icon="settings"
                        name="Settings"
                    />
                    <SidebarItem
                        onClick={() => setShowDeletePopup(true)}
                        icon="delete"
                        name="Delete"
                        red
                    />
                </SidebarGroup>
            </Sidebar>
            <div className={styles.side}>
                <Outlet />
            </div>
            <Popup show={showDeletePopup} onDismiss={dismissDeletePopup}>
                <Title>
                    Delete {instance?.display_name ?? instance?.service?.name}?
                </Title>
                <Text>
                    Are you sure you want to delete{" "}
                    {instance?.display_name ?? instance?.service?.name}? All
                    data will be permanently deleted.
                </Text>
                {isDeleting && <Progress infinite />}
                <APIError style={{ margin: 0 }} error={errorDeleting} />
                <Horizontal gap={10}>
                    <Spacer />
                    <Button onClick={dismissDeletePopup} disabled={isDeleting}>
                        Cancel
                    </Button>
                    <Button
                        primary
                        color="red"
                        onClick={async () => mutationDeleteInstance.mutate()}
                        disabled={isDeleting}
                        rightIcon="delete"
                    >
                        Confirm
                    </Button>
                </Horizontal>
            </Popup>
        </Horizontal>
    );

    return (
        <div className={styles.details}>
            <ProgressOverlay show={isLoading} />
            <div className={styles.bay}>
                <Bay
                    instances={[
                        {
                            value: instance,
                            onPower: async () => mutationInstancePower.mutate(),
                        },
                    ]}
                />
            </div>
            {!isLoading && content}
        </div>
    );
}
