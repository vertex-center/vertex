import { useState } from "react";
import { Outlet, useNavigate, useOutlet, useParams } from "react-router-dom";
import styles from "./Container.module.sass";
import { Horizontal } from "../../../../components/Layouts/Layouts";
import Popup, { PopupActions } from "../../../../components/Popup/Popup";
import {
    Button,
    MaterialIcon,
    Paragraph,
    Sidebar,
} from "@vertex-center/components";
import l from "../../../../components/NavLink/navlink";
import Progress from "../../../../components/Progress";
import useContainer from "../../hooks/useContainer";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useServerEvent } from "../../../../hooks/useEvent";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Container as ContainerModel } from "../../backend/models";
import Container from "../../../../components/Container/Container";
import { useSidebar } from "../../../../hooks/useSidebar";
import { API } from "../../backend/api";
import {
    ClockClockwise,
    Cube,
    Gear,
    House,
    ShareNetwork,
    TerminalWindow,
    Textbox,
    Trash,
} from "@phosphor-icons/react";

export default function ContainerDetails() {
    const { uuid } = useParams();
    const outlet = useOutlet();
    const navigate = useNavigate();
    const queryClient = useQueryClient();

    const { container, isLoading } = useContainer(uuid);

    const [showDeletePopup, setShowDeletePopup] = useState<boolean>();

    const route = uuid ? `/containers/${uuid}/events` : "";

    // @ts-ignore
    useServerEvent(window.api_urls.containers, route, {
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
                await API.startContainer(uuid);
            } else {
                await API.stopContainer(uuid);
            }
        },
    });

    const mutationDeleteContainer = useMutation({
        mutationFn: () => API.deleteContainer(uuid),
        onSuccess: () => navigate("/containers"),
    });
    const {
        isPending: isDeleting,
        error: errorDeleting,
        reset: resetDelete,
    } = mutationDeleteContainer;

    const dismissDeletePopup = () => {
        setShowDeletePopup(false);
        resetDelete();
    };

    const sidebar = useSidebar(
        <Sidebar>
            <Sidebar.Group>
                <Sidebar.Item
                    label="Home"
                    icon={<House />}
                    link={l(`/containers/${uuid}/home`)}
                />
            </Sidebar.Group>
            <Sidebar.Group title="Analyze">
                <Sidebar.Item
                    label="Logs"
                    icon={<TerminalWindow />}
                    link={l(`/containers/${uuid}/logs`)}
                />
                <Sidebar.Item
                    label="Docker"
                    icon={<Cube />}
                    link={l(`/containers/${uuid}/docker`)}
                />
            </Sidebar.Group>
            <Sidebar.Group title="Manage">
                <Sidebar.Item
                    label="Environment"
                    icon={<Textbox />}
                    link={l(`/containers/${uuid}/environment`)}
                />
                <Sidebar.Item
                    label="Ports"
                    icon={<ShareNetwork />}
                    link={l(`/containers/${uuid}/ports`)}
                />
                {container?.databases && (
                    <Sidebar.Item
                        label="Database"
                        icon={<MaterialIcon icon="database" />}
                        link={l(`/containers/${uuid}/database`)}
                    />
                )}
                <Sidebar.Item
                    icon={<ClockClockwise />}
                    label="Update"
                    link={l(`/containers/${uuid}/update`)}
                    notifications={
                        container?.service_update?.available ? 1 : undefined
                    }
                />
                <Sidebar.Item
                    label="Settings"
                    icon={<Gear />}
                    link={l(`/containers/${uuid}/settings`)}
                />
                <Sidebar.Item
                    label="Delete"
                    icon={<Trash />}
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
            <Popup
                show={showDeletePopup}
                onDismiss={dismissDeletePopup}
                title={`Delete ${container?.name}?`}
            >
                <Paragraph>
                    Are you sure you want to delete {container?.name}? All data
                    will be permanently deleted.
                </Paragraph>
                {isDeleting && <Progress infinite />}
                <APIError style={{ margin: 0 }} error={errorDeleting} />
                <PopupActions>
                    <Button onClick={dismissDeletePopup} disabled={isDeleting}>
                        Cancel
                    </Button>
                    <Button
                        variant="danger"
                        onClick={async () => mutationDeleteContainer.mutate()}
                        disabled={isDeleting}
                        rightIcon={<Trash />}
                    >
                        Confirm
                    </Button>
                </PopupActions>
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
