import { Title } from "../../components/Text/Text";
import { useParams } from "react-router-dom";
import useInstance from "../../hooks/useInstance";
import styles from "./BayDetailsUpdate.module.sass";
import { Vertical } from "../../components/Layouts/Layouts";
import { api } from "../../backend/backend";
import Update, { Updates } from "../../components/Update/Update";

export default function BayDetailsUpdate() {
    const { uuid } = useParams();

    const { instance, reloadInstance } = useInstance(uuid);

    const updateVertexIntegration = () => {
        return api.instance.update
            .service(uuid)
            .then(reloadInstance)
            .catch(console.error);
    };

    return (
        <Vertical gap={20}>
            <Title className={styles.title}>Update</Title>
            <Updates>
                <Update
                    name="Vertex integration"
                    onUpdate={updateVertexIntegration}
                    available={instance?.service_update?.available}
                />
            </Updates>
        </Vertical>
    );
}
