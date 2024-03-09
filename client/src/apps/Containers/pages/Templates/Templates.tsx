import { Template as ServiceModel } from "../../backend/template";
import React, { useState } from "react";
import styles from "./Templates.module.sass";
import Service from "../../../../components/Service/Service";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import ServiceInstallPopup from "../../../../components/ServiceInstallPopup/ServiceInstallPopup";
import { useQuery } from "@tanstack/react-query";
import {
    Button,
    Grid,
    Input,
    useTitle,
    Vertical,
} from "@vertex-center/components";
import { API } from "../../backend/api";
import Content from "../../../../components/Content/Content";
import { SiDocker } from "@icons-pack/react-simple-icons";
import ManualInstallPopup from "./ManualInstallPopup";
import { useCreateContainer } from "../../hooks/useCreateContainer";
import Spacer from "../../../../components/Spacer/Spacer";
import Toolbar from "../../../../components/Toolbar/Toolbar";

type Downloading = {
    service: ServiceModel;
};

export default function Templates() {
    useTitle("Create container");

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
            <Vertical gap={12} className={styles.content}>
                <APIError error={error} />
                <Toolbar>
                    <Input
                        placeholder="Search templates..."
                        style={{ width: 300 }}
                        disabled
                    />
                    <Spacer />
                    <Button
                        variant="colored"
                        rightIcon={<SiDocker size={20} />}
                        onClick={openManualInstallPopup}
                    >
                        Or install manually
                    </Button>
                </Toolbar>
                <Grid rowSize={200}>
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
                </Grid>
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
