import { Template as ServiceModel } from "../../backend/template";
import React, { useState } from "react";
import styles from "./Templates.module.sass";
import Service from "../../../../components/Service/Service";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import TemplateInstallPopup from "../../../../components/TemplateInstallPopup/TemplateInstallPopup";
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
import Spacer from "../../../../components/Spacer/Spacer";
import Toolbar from "../../../../components/Toolbar/Toolbar";

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

    const [showInstallPopup, setShowInstallPopup] = useState<boolean>(false);
    const [showManualInstallPopup, setShowManualInstallPopup] =
        useState<boolean>(false);
    const [selectedTemplate, setSelectedTemplate] = useState<ServiceModel>();

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

    const error = servicesError;

    return (
        <Content>
            <ProgressOverlay show={isServicesLoading} />
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
                <Grid columnSize={200}>
                    {services?.map((template) => (
                        <Service
                            key={template.id}
                            template={template}
                            onInstall={() => openInstallPopup(template)}
                        />
                    ))}
                </Grid>
            </Vertical>
            {showInstallPopup && (
                <TemplateInstallPopup
                    service={selectedTemplate}
                    dismiss={closeInstallPopup}
                />
            )}
            {showManualInstallPopup && (
                <ManualInstallPopup dismiss={closeManualInstallPopup} />
            )}
        </Content>
    );
}
