import { Service as ServiceModel } from "../../backend/service";
import React, { useState } from "react";
import styles from "./ContainersStore.module.sass";
import Service from "../../../../components/Service/Service";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import ServiceInstallPopup from "../../../../components/ServiceInstallPopup/ServiceInstallPopup";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
    List,
    ListActions,
    ListDescription,
    ListIcon,
    ListInfo,
    ListItem,
    ListTitle,
    MaterialIcon,
    Title,
    useTitle,
    Vertical,
} from "@vertex-center/components";
import { API } from "../../backend/api";
import Content from "../../../../components/Content/Content";
import { SiDocker } from "@icons-pack/react-simple-icons";
import ManualInstallPopup from "./ManualInstallPopup";

type Downloading = {
    service: ServiceModel;
};

export default function ContainersStore() {
    useTitle("Create container");

    const queryClient = useQueryClient();

    const queryServices = useQuery({
        queryKey: ["services"],
        queryFn: API.getAllServices,
    });
    const {
        data: services,
        isLoading: isServicesLoading,
        error: servicesError,
    } = queryServices;

    const queryContainers = useQuery({
        queryKey: ["containers"],
        queryFn: () => API.getContainers(),
    });
    const {
        data: containers,
        isLoading: isContainersLoading,
        error: containersError,
    } = queryContainers;

    const mutationCreateContainer = useMutation({
        mutationFn: API.createContainer,
        onSettled: (data, error, serviceID) => {
            setDownloading(
                downloading.filter(({ service: s }) => s.id !== serviceID)
            );
            queryClient.invalidateQueries({
                queryKey: ["containers"],
            });
        },
    });
    const { isLoading: isInstalling, error: installError } =
        mutationCreateContainer;

    const install = () => {
        const service = selectedService;
        setDownloading((prev) => [...prev, { service }]);
        setShowInstallPopup(false);
        mutationCreateContainer.mutate(service.id);
    };

    const [showInstallPopup, setShowInstallPopup] = useState<boolean>(false);
    const [showManualInstallPopup, setShowManualInstallPopup] =
        useState<boolean>(false);
    const [selectedService, setSelectedService] = useState<ServiceModel>();
    const [downloading, setDownloading] = useState<Downloading[]>([]);

    const openInstallPopup = (service: ServiceModel) => {
        setSelectedService(service);
        setShowInstallPopup(true);
    };

    const closeInstallPopup = () => {
        setSelectedService(undefined);
        setShowInstallPopup(false);
    };

    const openManualInstallPopup = () => {
        setShowManualInstallPopup(true);
    };

    const closeManualInstallPopup = () => {
        setShowManualInstallPopup(false);
    };

    const error = servicesError ?? containersError ?? installError;

    return (
        <Content>
            <ProgressOverlay
                show={isContainersLoading ?? isServicesLoading ?? isInstalling}
            />
            <Vertical gap={30} className={styles.content}>
                <List>
                    <ListItem onClick={openManualInstallPopup}>
                        <ListIcon>
                            <SiDocker />
                        </ListIcon>
                        <ListInfo>
                            <ListTitle>
                                Manually from a Docker Registry
                            </ListTitle>
                            <ListDescription>
                                This will need manual configuration
                            </ListDescription>
                        </ListInfo>
                        <ListActions>
                            <MaterialIcon icon="download" />
                        </ListActions>
                    </ListItem>
                </List>
                <Title variant="h2">From template</Title>
                <APIError error={error} />
                <List>
                    {services?.map((serv) => (
                        <Service
                            key={serv.id}
                            service={serv}
                            onInstall={() => openInstallPopup(serv)}
                            downloading={downloading.some(
                                ({ service: s }) => s.id === serv.id
                            )}
                            installedCount={
                                containers?.filter(
                                    (c) => c.service_id === serv.id
                                )?.length
                            }
                        />
                    ))}
                </List>
            </Vertical>
            <ServiceInstallPopup
                service={selectedService}
                show={showInstallPopup}
                dismiss={closeInstallPopup}
                install={install}
            />
            <ManualInstallPopup
                show={showManualInstallPopup}
                dismiss={closeManualInstallPopup}
            />
        </Content>
    );
}
