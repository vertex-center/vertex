import { Service as ServiceModel } from "../../backend/service";
import { Fragment, useState } from "react";
import styles from "./ContainersStore.module.sass";
import Service from "../../../../components/Service/Service";
import { Vertical } from "../../../../components/Layouts/Layouts";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import ServiceInstallPopup from "../../../../components/ServiceInstallPopup/ServiceInstallPopup";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { List, useTitle } from "@vertex-center/components";
import { API } from "../../backend/api";

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
        queryFn: API.getAllContainers,
    });
    const {
        data: containers,
        isLoading: isContainersLoading,
        error: containersError,
    } = queryContainers;

    const mutationInstallService = useMutation({
        mutationFn: API.installService,
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

    const error = servicesError ?? containersError ?? installError;

    return (
        <Fragment>
            <ProgressOverlay
                show={isContainersLoading ?? isServicesLoading ?? isInstalling}
            />
            <Vertical className={styles.content} gap={10}>
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
        </Fragment>
    );
}
