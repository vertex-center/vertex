import { Title } from "../../../components/Text/Text";
import { api } from "../../../backend/backend";
import { useFetch } from "../../../hooks/useFetch";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import Service from "../../../components/Service/Service";
import { Service as ServiceModel } from "../../../models/service";
import styles from "./SqlInstaller.module.sass";
import { Vertical } from "../../../components/Layouts/Layouts";
import ServiceInstallPopup from "../../../components/ServiceInstallPopup/ServiceInstallPopup";
import { useState } from "react";
import { Instances } from "../../../models/instance";
import { APIError } from "../../../components/Error/APIError";

export default function SqlInstaller() {
    const {
        data: services,
        loading,
        error: servicesError,
    } = useFetch<ServiceModel[]>(api.services.available.get);
    const {
        data: instances,
        error: instancesError,
        reload: reloadInstances,
    } = useFetch<Instances>(api.instances.get);

    const [selectedService, setSelectedService] = useState<ServiceModel>();
    const [showPopup, setShowPopup] = useState(false);

    const [error, setError] = useState();
    const [downloading, setDownloading] = useState<{ service: ServiceModel }[]>(
        []
    );

    const open = (service: ServiceModel) => {
        setSelectedService(service);
        setShowPopup(true);
    };

    const dismiss = () => {
        setSelectedService(undefined);
        setShowPopup(false);
    };

    const install = () => {
        const service = selectedService;

        setDownloading((prev) => [...prev, { service }]);
        setShowPopup(false);

        api.services
            .install({
                method: "docker",
                service_id: service.id,
            })
            .catch(setError)
            .finally(() => {
                setDownloading((d) =>
                    d?.filter(({ service: s }) => s.id !== service.id)
                );
            });
    };

    return (
        <Vertical gap={20}>
            <ProgressOverlay show={loading} />
            <Title className={styles.title}>SQL Database Installer</Title>
            <APIError error={error ?? servicesError ?? instancesError} />
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
