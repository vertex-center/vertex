import { api } from "../../../backend/api/backend";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import Service from "../../../components/Service/Service";
import { Service as ServiceModel } from "../../Containers/backend/service";
import ServiceInstallPopup from "../../../components/ServiceInstallPopup/ServiceInstallPopup";
import { useState } from "react";
import { APIError } from "../../../components/Error/APIError";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { List, Title } from "@vertex-center/components";
import Content from "../../../components/Content/Content";

export default function SqlInstaller() {
    const queryClient = useQueryClient();

    const queryServices = useQuery({
        queryKey: ["services"],
        queryFn: api.vxContainers.services.all,
    });
    const {
        data: services,
        isLoading: isServicesLoading,
        error: servicesError,
    } = queryServices;

    const queryContainers = useQuery({
        queryKey: ["containers"],
        queryFn: api.vxContainers.containers.all,
    });
    const {
        data: containers,
        isLoading: isContainersLoading,
        error: containersError,
    } = queryContainers;

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

    const mutationInstallService = useMutation({
        mutationFn: async (serviceId: string) => {
            await api.vxSql.dbms(serviceId).install();
        },
        onSettled: (data, error, serviceId) => {
            setDownloading(
                downloading.filter(({ service: s }) => s.id !== serviceId)
            );
            queryClient.invalidateQueries({
                queryKey: ["containers"],
            });
        },
    });
    const { isLoading: isInstalling, error: installError } =
        mutationInstallService;

    const install = () => {
        const service = selectedService;
        setDownloading((prev) => [...prev, { service }]);
        setShowPopup(false);
        mutationInstallService.mutate(service.id);
    };

    const error = servicesError ?? containersError ?? installError;

    return (
        <Content>
            <ProgressOverlay show={isContainersLoading ?? isServicesLoading} />
            <Title variant="h2">Installer</Title>
            <APIError error={error} />
            <List>
                {services
                    ?.filter((s) => s?.features?.databases?.length >= 1)
                    ?.filter((s) =>
                        s?.features?.databases?.some(
                            (d) => d.category === "sql"
                        )
                    )
                    ?.map((service) => {
                        return (
                            <Service
                                key={service.id}
                                service={service}
                                onInstall={() => open(service)}
                                downloading={downloading.some(
                                    ({ service: s }) => s.id === service.id
                                )}
                                installedCount={
                                    containers === undefined
                                        ? undefined
                                        : Object.values(containers)?.filter(
                                              ({ service: s }) =>
                                                  s.id === service.id
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
