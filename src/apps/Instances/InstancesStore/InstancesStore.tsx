import { BigTitle } from "../../../components/Text/Text";
import { Service as ServiceModel } from "../../../models/service";
import { Fragment, useState } from "react";
import styles from "./InstancesStore.module.sass";
import Service from "../../../components/Service/Service";
import { Horizontal, Vertical } from "../../../components/Layouts/Layouts";
import { api } from "../../../backend/backend";
import { APIError } from "../../../components/Error/APIError";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import ServiceInstallPopup from "../../../components/ServiceInstallPopup/ServiceInstallPopup";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

type Downloading = {
    service: ServiceModel;
};

export default function InstancesStore() {
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

    const mutationInstallService = useMutation({
        mutationFn: async (serviceId: string) => {
            await api.vxInstances.service(serviceId).install();
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
        setShowInstallPopup(false);
        mutationInstallService.mutate(service.id);
    };

    const [showInstallPopup, setShowInstallPopup] = useState<boolean>(false);
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

    const error = servicesError ?? instancesError ?? installError;

    return (
        <Fragment>
            <ProgressOverlay
                show={isInstancesLoading ?? isServicesLoading ?? isInstalling}
            />
            <div className={styles.page}>
                <Horizontal
                    className={styles.title}
                    gap={10}
                    alignItems="center"
                >
                    <BigTitle>Create instance</BigTitle>
                </Horizontal>
                <Vertical className={styles.content}>
                    <APIError error={error} />
                    {services?.map((service) => (
                        <Service
                            key={service.id}
                            service={service}
                            onInstall={() => openInstallPopup(service)}
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
                    ))}
                </Vertical>
            </div>
            <ServiceInstallPopup
                service={selectedService}
                show={showInstallPopup}
                dismiss={closeInstallPopup}
                install={install}
            />
        </Fragment>
    );
}
