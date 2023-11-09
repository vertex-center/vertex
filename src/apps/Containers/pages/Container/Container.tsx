import { useState } from "react";
import { api } from "../../../../backend/api/backend";
import { Outlet, useNavigate, useOutlet, useParams } from "react-router-dom";
import styles from "./Container.module.sass";
import { Horizontal } from "../../../../components/Layouts/Layouts";
import Spacer from "../../../../components/Spacer/Spacer";
import Popup from "../../../../components/Popup/Popup";
import { Text, Title } from "../../../../components/Text/Text";
import { Button, MaterialIcon, Sidebar } from "@vertex-center/components";
import l from "../../../../components/NavLink/navlink";
import Progress from "../../../../components/Progress";
import { SiDocker } from "@icons-pack/react-simple-icons";
import useContainer from "../../hooks/useContainer";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useServerEvent } from "../../../../hooks/useEvent";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Container as ContainerModel } from "../../../../models/container";
import Container from "../../../../components/Container/Container";
import { useSidebar } from "../../../../hooks/useSidebar";

export default function ContainerDetails() {
    const { uuid } = useParams();
    const outlet = useOutlet();
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

    const sidebar = useSidebar(
        <Sidebar>
            <Sidebar.Group>
                <Sidebar.Item
                    label="Home"
                    icon={<MaterialIcon icon="home" />}
                    link={l(`/app/vx-containers/${uuid}/home`)}
                />
            </Sidebar.Group>
            <Sidebar.Group title="Analyze">
                <Sidebar.Item
                    label="Logs"
                    icon={<MaterialIcon icon="terminal" />}
                    link={l(`/app/vx-containers/${uuid}/logs`)}
                />
                {container?.install_method === "docker" && (
                    <Sidebar.Item
                        label="Docker"
                        icon={<SiDocker size={20} />}
                        link={l(`/app/vx-containers/${uuid}/docker`)}
                    />
                )}
            </Sidebar.Group>
            <Sidebar.Group title="Manage">
                <Sidebar.Item
                    label="Environment"
                    icon={<MaterialIcon icon="tune" />}
                    link={l(`/app/vx-containers/${uuid}/environment`)}
                />
                {container?.service?.databases && (
                    <Sidebar.Item
                        label="Database"
                        icon={<MaterialIcon icon="database" />}
                        link={l(`/app/vx-containers/${uuid}/database`)}
                    />
                )}
                <Sidebar.Item
                    icon={<MaterialIcon icon="update" />}
                    label="Update"
                    link={l(`/app/vx-containers/${uuid}/update`)}
                    notifications={
                        container?.service_update?.available ? 1 : undefined
                    }
                />
                <Sidebar.Item
                    label="Settings"
                    icon={<MaterialIcon icon="settings" />}
                    link={l(`/app/vx-containers/${uuid}/settings`)}
                />
                <Sidebar.Item
                    label="Delete"
                    icon={<MaterialIcon icon="delete" />}
                    onClick={() => setShowDeletePopup(true)}
                    variant="red"
                />
            </Sidebar.Group>
        </Sidebar>
    );

    const content = (
        <Horizontal className={styles.content}>
            {sidebar}
            {outlet && (
                <div className={styles.side}>
                    <Outlet />
                </div>
            )}
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
                        variant="danger"
                        onClick={async () => mutationDeleteContainer.mutate()}
                        disabled={isDeleting}
                        rightIcon={<MaterialIcon icon="delete" />}
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
