import { ProgressOverlay } from "../../../components/Progress/Progress";
import Service from "../../../components/Service/Service";
import { Template as ServiceModel } from "../../Containers/backend/template";
import TemplateInstallPopup from "../../../components/TemplateInstallPopup/TemplateInstallPopup";
import { useState } from "react";
import { APIError } from "../../../components/Error/APIError";
import { useQuery } from "@tanstack/react-query";
import { List, Title } from "@vertex-center/components";
import Content from "../../../components/Content/Content";
import { API } from "../../Containers/backend/api";

export default function SqlInstaller() {
    const queryServices = useQuery({
        queryKey: ["templates"],
        queryFn: API.getAllTemplates,
    });
    const {
        data: templates,
        isLoading: isServicesLoading,
        error: servicesError,
    } = queryServices;

    const [selectedService, setSelectedService] = useState<ServiceModel>();
    const [showPopup, setShowPopup] = useState(false);

    const open = (service: ServiceModel) => {
        setSelectedService(service);
        setShowPopup(true);
    };

    const dismiss = () => {
        setSelectedService(undefined);
        setShowPopup(false);
    };

    return (
        <Content>
            <ProgressOverlay show={isServicesLoading} />
            <Title variant="h2">Installer</Title>
            <APIError error={servicesError} />
            <List>
                {templates
                    ?.filter((s) => s?.features?.databases?.length >= 1)
                    ?.filter((s) =>
                        s?.features?.databases?.some(
                            (d) => d.category === "sql"
                        )
                    )
                    ?.map((template) => {
                        return (
                            <Service
                                key={template.id}
                                template={template}
                                onInstall={() => open(template)}
                            />
                        );
                    })}
            </List>
            {showPopup && selectedService && (
                <TemplateInstallPopup
                    service={selectedService}
                    dismiss={dismiss}
                />
            )}
        </Content>
    );
}
