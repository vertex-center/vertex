import { Service as ServiceModel } from "../../../../models/service";
import { Fragment, useState } from "react";
import styles from "./ContainersStore.module.sass";
import Service from "../../../../components/Service/Service";
import { Vertical } from "../../../../components/Layouts/Layouts";
import { api } from "../../../../backend/api/backend";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import ServiceInstallPopup from "../../../../components/ServiceInstallPopup/ServiceInstallPopup";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import List from "../../../../components/List/List";
import { useTitle } from "@vertex-center/components";

type Downloading = {
    service: ServiceModel;
};

export default function ContainersStore() {
    useTitle("Create container");

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

    const mutationInstallService = useMutation({
        mutationFn: async (serviceId: string) => {
            await api.vxContainers.service(serviceId).install();
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
            <div className={styles.page}>
                <Vertical gap={10}>
                    <APIError error={error} />
                    {/*<Toolbar className={styles.toolbar}>*/}
                    {/*    <Select*/}
                    {/*        // @ts-ignore*/}
                    {/*        value={<SelectValue>Category</SelectValue>}*/}
                    {/*        disabled*/}
                    {/*    />*/}
                    {/*    <Spacer />*/}
                    {/*    <Button*/}
                    {/*        to="/app/vx-devtools-service-editor"*/}
                    {/*        leftIcon="frame_source"*/}
                    {/*    >*/}
                    {/*        Service Editor*/}
                    {/*    </Button>*/}
                    {/*</Toolbar>*/}
                    <List className={styles.content}>
                        {services?.map((service) => (
                            <Service
                                key={service.id}
                                service={service}
                                onInstall={() => openInstallPopup(service)}
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
                        ))}
                    </List>
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
