import { Title } from "../../../../components/Text/Text";
import { useParams } from "react-router-dom";
import useInstance from "../../../../hooks/useInstance";
import styles from "./InstanceUpdate.module.sass";
import { Vertical } from "../../../../components/Layouts/Layouts";
import { api } from "../../../../backend/backend";
import Update, { Updates } from "../../../../components/Update/Update";
import { useState } from "react";
import { APIError } from "../../../../components/Error/Error";
import { ProgressOverlay } from "../../../../components/Progress/Progress";

export default function InstanceUpdate() {
    const { uuid } = useParams();

    const { instance, reloadInstance, loading } = useInstance(uuid);

    const [error, setError] = useState();

    const updateVertexIntegration = () => {
        return api.instance.update
            .service(uuid)
            .then(reloadInstance)
            .catch(setError);
    };

    return (
        <Vertical gap={20}>
            <ProgressOverlay show={loading} />
            <Title className={styles.title}>Update</Title>
            <APIError error={error} />
            {!error && !loading && (
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
