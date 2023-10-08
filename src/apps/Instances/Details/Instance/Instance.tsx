import Bay from "../../../../components/Bay/Bay";
import { useState } from "react";
import { api } from "../../../../backend/backend";
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

export default function Instance() {
    const { uuid } = useParams();
    const navigate = useNavigate();

    const { instance, setInstance, loading } = useInstance(uuid);

    const [showDeletePopup, setShowDeletePopup] = useState<boolean>();
    const [deleting, setDeleting] = useState<boolean>(false);
    const [error, setError] = useState();

    const route = uuid ? `/instance/${uuid}/events` : "";

    useServerEvent(route, {
        status_change: (e) => {
            setInstance((instance) => ({ ...instance, status: e.data }));
        },
    });

    const toggleInstance = async (uuid: string) => {
        if (instance.status === "off" || instance.status === "error") {
            await api.instance.start(uuid);
        } else {
            await api.instance.stop(uuid);
        }
    };

    const onDeleteInstance = () => {
        setDeleting(true);
        setError(undefined);
        api.instance
            .delete(uuid)
            .then(() => {
                navigate("/instances");
            })
            .catch((error) => {
                setError(error);
                setDeleting(false);
            });
    };

    const dismissDeletePopup = () => {
        setShowDeletePopup(false);
        setError(undefined);
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
                {deleting && <Progress infinite />}
                <APIError style={{ margin: 0 }} error={error} />
                <Horizontal gap={10}>
                    <Spacer />
                    <Button onClick={dismissDeletePopup} disabled={deleting}>
                        Cancel
                    </Button>
                    <Button
                        primary
                        color="red"
                        onClick={onDeleteInstance}
                        disabled={deleting}
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
            <ProgressOverlay show={loading} />
            <div className={styles.bay}>
                <Bay
                    instances={[
                        {
                            value: instance,
                            onPower: () => toggleInstance(uuid),
                        },
                    ]}
                />
            </div>
            {!loading && content}
        </div>
    );
}
