import Container from "../../../../components/Container/Container";
import { useState } from "react";
import { api } from "../../../../backend/api/backend";
import { Outlet, useNavigate, useParams } from "react-router-dom";
import styles from "./Container.module.sass";
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
import useContainer from "../../../../hooks/useContainer";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useServerEvent } from "../../../../hooks/useEvent";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Container as ContainerModel } from "../../../../models/container";

export default function ContainerDetails() {
    const { uuid } = useParams();
    const navigate = useNavigate();
    const queryClient = useQueryClient();

    const { container, isLoading } = useContainer(uuid);

    const [showDeletePopup, setShowDeletePopup] = useState<boolean>();

    const route = uuid ? `/app/vx-containers/container/${uuid}/events` : "";

    useServerEvent(route, {
        status_change: (e) => {
            queryClient.setQueryData(
                ["containers", uuid],
                (container: ContainerModel) => ({
                    ...container,
                    status: e.data,
                })
            );
        },
    });

    const mutationContainerPower = useMutation({
        mutationFn: async () => {
            if (container.status === "off" || container.status === "error") {
                await api.vxContainers.container(uuid).start();
            } else {
                await api.vxContainers.container(uuid).stop();
            }
        },
    });

    const mutationDeleteContainer = useMutation({
        mutationFn: api.vxContainers.container(uuid).delete,
        onSuccess: () => {
            navigate("/app/vx-containers");
        },
    });
    const { isLoading: isDeleting, error: errorDeleting } =
        mutationDeleteContainer;

    const dismissDeletePopup = () => {
        setShowDeletePopup(false);
    };

    const content = (
        <Horizontal className={styles.content}>
            <Sidebar root={`/app/vx-containers/${uuid}`}>
                <SidebarGroup>
                    <SidebarItem
                        to={`/app/vx-containers/${uuid}/home`}
                        icon="home"
                        name="Home"
                    />
                </SidebarGroup>
                <SidebarGroup title="Analyze">
                    <SidebarItem
                        to={`/app/vx-containers/${uuid}/logs`}
                        icon="terminal"
                        name="Logs"
                    />
                    {container?.install_method === "docker" && (
                        <SidebarItem
                            to={`/app/vx-containers/${uuid}/docker`}
                            icon={<SiDocker size={20} />}
                            name="Docker"
                        />
                    )}
                </SidebarGroup>
                <SidebarGroup title="Manage">
                    <SidebarItem
                        to={`/app/vx-containers/${uuid}/environment`}
                        icon="tune"
                        name="Environment"
                    />
                    {container?.service?.databases && (
                        <SidebarItem
                            to={`/app/vx-containers/${uuid}/database`}
                            icon="database"
                            name="Database"
                        />
                    )}
                    <SidebarItem
                        to={`/app/vx-containers/${uuid}/update`}
                        icon="update"
                        name="Update"
                        notifications={
                            container?.service_update?.available ? 1 : undefined
                        }
                    />
                    <SidebarItem
                        to={`/app/vx-containers/${uuid}/settings`}
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
                    Delete {container?.display_name ?? container?.service?.name}
                    ?
                </Title>
                <Text>
                    Are you sure you want to delete{" "}
                    {container?.display_name ?? container?.service?.name}? All
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
                        onClick={async () => mutationDeleteContainer.mutate()}
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
            <div className={styles.container}>
                <Container
                    container={{
                        value: container,
                        onPower: async () => mutationContainerPower.mutate(),
                    }}
                />
            </div>
            {!isLoading && content}
        </div>
    );
}
