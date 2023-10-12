import { Title } from "../../../../components/Text/Text";
import { useParams } from "react-router-dom";
import useInstance from "../../../../hooks/useInstance";
import styles from "./InstanceUpdate.module.sass";
import { Vertical } from "../../../../components/Layouts/Layouts";
import { api } from "../../../../backend/api/backend";
import Update, { Updates } from "../../../../components/Update/Update";
import { useState } from "react";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useQueryClient } from "@tanstack/react-query";

export default function InstanceUpdate() {
    const { uuid } = useParams();
    const queryClient = useQueryClient();

    const { instance, isLoading } = useInstance(uuid);

    const [error, setError] = useState();

    const updateVertexIntegration = () => {
        return api.vxInstances
            .instance(uuid)
            .update.service()
            .then(() => {
                queryClient.invalidateQueries({
                    queryKey: ["instances", uuid],
                });
            })
            .catch(setError);
    };

    return (
        <Vertical gap={20}>
            <ProgressOverlay show={isLoading} />
            <Title className={styles.title}>Update</Title>
            <APIError error={error} />
            {!error && !isLoading && (
                <Updates>
                    <Update
                        name="Vertex integration"
                        onUpdate={updateVertexIntegration}
                        available={instance?.service_update?.available}
                    />
                </Updates>
            )}
        </Vertical>
    );
}
