import { Template as ServiceModel } from "../../backend/template";
import React, { useState } from "react";
import styles from "./ContainersStore.module.sass";
import Service from "../../../../components/Service/Service";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import ServiceInstallPopup from "../../../../components/ServiceInstallPopup/ServiceInstallPopup";
import { useQuery, useQueryClient } from "@tanstack/react-query";
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
import { useCreateContainer } from "../../hooks/useCreateContainer";

type Downloading = {
    service: ServiceModel;
};

export default function ContainersStore() {
    useTitle("Create container");

    const queryClient = useQueryClient();

    const queryServices = useQuery({
        queryKey: ["services"],
        queryFn: API.getAllTemplates,
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

    const { createContainer, isCreatingContainer, errorCreatingContainer } =
        useCreateContainer({
            onSettled: (data, error, options) => {
                setDownloading(
                    downloading.filter(
                        ({ service: s }) => s.id !== options.template_id
                    )
                );
            },
        });

    const install = () => {
        const template = selectedTemplate;
        setDownloading((prev) => [...prev, { service: template }]);
        setShowInstallPopup(false);
        createContainer({
            template_id: template.id,
        });
    };

    const [showInstallPopup, setShowInstallPopup] = useState<boolean>(false);
    const [showManualInstallPopup, setShowManualInstallPopup] =
        useState<boolean>(false);
    const [selectedTemplate, setSelectedTemplate] = useState<ServiceModel>();
    const [downloading, setDownloading] = useState<Downloading[]>([]);

    const openInstallPopup = (template: ServiceModel) => {
        setSelectedTemplate(template);
        setShowInstallPopup(true);
    };

    const closeInstallPopup = () => {
        setSelectedTemplate(undefined);
        setShowInstallPopup(false);
    };

    const openManualInstallPopup = () => {
        setShowManualInstallPopup(true);
    };

    const closeManualInstallPopup = () => {
        setShowManualInstallPopup(false);
    };

    const error = servicesError ?? containersError ?? errorCreatingContainer;

    return (
        <Content>
            <ProgressOverlay
                show={
                    isContainersLoading ??
                    isServicesLoading ??
                    isCreatingContainer
                }
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
                    {services?.map((template) => (
                        <Service
                            key={template.id}
                            template={template}
                            onInstall={() => openInstallPopup(template)}
                            downloading={downloading.some(
                                ({ service: s }) => s.id === template.id
                            )}
                            installedCount={
                                containers?.filter(
                                    (c) => c.template_id === template.id
                                )?.length
                            }
                        />
                    ))}
                </List>
            </Vertical>
            <ServiceInstallPopup
                service={selectedTemplate}
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
