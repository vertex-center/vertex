import { Title } from "../../../components/Text/Text";
import { api } from "../../../backend/backend";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import Service from "../../../components/Service/Service";
import { Service as ServiceModel } from "../../../models/service";
import styles from "./SqlInstaller.module.sass";
import { Vertical } from "../../../components/Layouts/Layouts";
import ServiceInstallPopup from "../../../components/ServiceInstallPopup/ServiceInstallPopup";
import { useState } from "react";
import { APIError } from "../../../components/Error/APIError";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export default function SqlInstaller() {
    const queryClient = useQueryClient();

    const queryServices = useQuery({
        queryKey: ["services"],
        queryFn: api.vxInstances.services.all,
    });
    const {
        data: services,
        isLoading: isServicesLoading,
        error: servicesError,
    } = queryServices;

    const queryInstances = useQuery({
        queryKey: ["instances"],
        queryFn: api.vxInstances.instances.all,
    });
    const {
        data: instances,
        isLoading: isInstancesLoading,
        error: instancesError,
    } = queryInstances;

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
                queryKey: ["instances"],
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

    const error = servicesError ?? instancesError ?? installError;

    return (
        <Vertical gap={20}>
            <ProgressOverlay show={isInstancesLoading ?? isServicesLoading} />
            <Title className={styles.title}>SQL Database Installer</Title>
            <APIError error={error} />
            <Vertical>
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
                                    instances === undefined
                                        ? undefined
                                        : Object.values(instances)?.filter(
                                              ({ service: s }) =>
                                                  s.id === service.id
                                          )?.length
                                }
                            />
                        );
                    })}
            </Vertical>
            <ServiceInstallPopup
                service={selectedService}
                show={showPopup}
                dismiss={dismiss}
                install={install}
            />
        </Vertical>
    );
}