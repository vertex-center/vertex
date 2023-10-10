import { useFetch } from "../../../hooks/useFetch";
import { BigTitle } from "../../../components/Text/Text";
import { Service as ServiceModel } from "../../../models/service";
import { Fragment, useState } from "react";

import styles from "./InstancesStore.module.sass";
import Service from "../../../components/Service/Service";
import { Horizontal, Vertical } from "../../../components/Layouts/Layouts";
import { Instances } from "../../../models/instance";
import { api } from "../../../backend/backend";
import { Errors } from "../../../components/Error/Errors";
import { APIError } from "../../../components/Error/APIError";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import ServiceInstallPopup from "../../../components/ServiceInstallPopup/ServiceInstallPopup";

type Downloading = {
    service: ServiceModel;
};

export default function InstancesStore() {
    const {
        data: services,
        error: servicesError,
        loading,
    } = useFetch<ServiceModel[]>(api.vxInstances.services.all);
    const {
        data: instances,
        error: instancesError,
        reload: reloadInstances,
    } = useFetch<Instances>(api.vxInstances.instances.all);

    const [showInstallPopup, setShowInstallPopup] = useState<boolean>(false);

    const [selectedService, setSelectedService] = useState<ServiceModel>();

    const [error, setError] = useState();

    const [downloading, setDownloading] = useState<Downloading[]>([]);

    const openInstallPopup = (service: ServiceModel) => {
        setSelectedService(service);
        setShowInstallPopup(true);
    };

    const closeInstallPopup = () => {
        setSelectedService(undefined);
        setShowInstallPopup(false);
    };

    const install = () => {
        const service = selectedService;

        setDownloading((prev) => [...prev, { service }]);
        setShowInstallPopup(false);

        api.vxInstances
            .service(service.id)
            .install()
            .catch(setError)
            .finally(() => {
                setDownloading(
                    downloading.filter(({ service: s }) => s.id !== service.id)
                );
                reloadInstances().catch(console.error);
            });
    };

    return (
        <Fragment>
            <ProgressOverlay show={loading} />
            <div className={styles.page}>
                <Horizontal
                    className={styles.title}
                    gap={10}
                    alignItems="center"
                >
                    <BigTitle>Create instance</BigTitle>
                </Horizontal>
                <Errors className={styles.errors}>
                    <APIError error={error} />
                    <APIError error={servicesError} />
                    <APIError error={instancesError} />
                </Errors>
                <Vertical className={styles.content}>
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
