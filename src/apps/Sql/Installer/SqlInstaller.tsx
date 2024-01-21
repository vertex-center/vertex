import { api } from "../../../backend/api/backend";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import Service from "../../../components/Service/Service";
import { Template as ServiceModel } from "../../Containers/backend/template";
import ServiceInstallPopup from "../../../components/ServiceInstallPopup/ServiceInstallPopup";
import { useState } from "react";
import { APIError } from "../../../components/Error/APIError";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { List, Title } from "@vertex-center/components";
import Content from "../../../components/Content/Content";
import { API } from "../../Containers/backend/api";
import { useContainers } from "../../Containers/hooks/useContainers";

export default function SqlInstaller() {
    const queryClient = useQueryClient();

    const queryServices = useQuery({
        queryKey: ["templates"],
        queryFn: API.getAllTemplates,
    });
    const {
        data: templates,
        isLoading: isServicesLoading,
        error: servicesError,
    } = queryServices;

    const {
        containers,
        isLoading: isLoadingContainers,
        error: containersError,
    } = useContainers({});

    const [selectedService, setSelectedService] = useState<ServiceModel>();
    const [showPopup, setShowPopup] = useState(false);
    const [downloading, setDownloading] = useState<
        {
            service: ServiceModel;
        }[]
    >([]);

    const open = (service: ServiceModel) => {
        setSelectedService(service);
        setShowPopup(true);
    };

    const dismiss = () => {
        setSelectedService(undefined);
        setShowPopup(false);
    };

    const mutationCreateContainer = useMutation({
        mutationFn: async (templateID: string) => {
            await api.vxSql.dbms(templateID).install();
        },
        onSettled: (data, error, templateID) => {
            setDownloading(
                downloading.filter(({ service: s }) => s.id !== templateID)
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
        setShowPopup(false);
        mutationCreateContainer.mutate(service.id);
    };

    const error = servicesError ?? containersError ?? installError;

    return (
        <Content>
            <ProgressOverlay show={isLoadingContainers ?? isServicesLoading} />
            <Title variant="h2">Installer</Title>
            <APIError error={error} />
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
                                downloading={downloading.some(
                                    ({ service: s }) => s.id === template.id
                                )}
                                installedCount={
                                    containers === undefined
                                        ? undefined
                                        : Object.values(
                                              containers ?? []
                                          )?.filter(
                                              ({ template_id }) =>
                                                  template_id === template.id
                                          )?.length
                                }
                            />
                        );
                    })}
            </List>
            <ServiceInstallPopup
                service={selectedService}
                show={showPopup}
                dismiss={dismiss}
                install={install}
            />
        </Content>
    );
}
