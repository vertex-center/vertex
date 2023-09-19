import { SubTitle, Title } from "../../components/Text/Text";
import { useParams } from "react-router-dom";
import useInstance from "../../hooks/useInstance";
import styles from "./BayDetailsUpdate.module.sass";
import Button from "../../components/Button/Button";
import { Vertical } from "../../components/Layouts/Layouts";
import { Fragment, useState } from "react";
import { api } from "../../backend/backend";

export default function BayDetailsUpdate() {
    const { uuid } = useParams();

    const { instance, reloadInstance } = useInstance(uuid);

    const [updating, setUpdating] = useState(false);

    // let content: any;
    // if (instance?.update) {
    //     content = (
    //         <Vertical gap={10}>
    //             <Text>An update is available.</Text>
    //             <Text>Current: {instance?.update?.current_version}</Text>
    //             <Text>Latest: {instance?.update?.latest_version}</Text>
    //         </Vertical>
    //     );
    // } else {
    //     content = (
    //         <Text className={styles.content}>Everything is up-to-date</Text>
    //     );
    // }

    const updateVertexIntegration = () => {
        setUpdating(true);
        api.instance.update
            .service(uuid)
            .then(reloadInstance)
            .catch(console.error)
            .finally(() => setUpdating(false));
    };

    return (
        <Vertical gap={20}>
            <Title className={styles.title}>Update</Title>
            {instance?.service_update?.available && (
                <Fragment>
                    <SubTitle className={styles.content}>
                        Vertex integration
                    </SubTitle>
                    <div className={styles.content}>
                        A new version of the Vertex integration is available for
                        this instance.
                    </div>
                    <div>
                        <Button
                            disabled={updating}
                            rightSymbol="download"
                            onClick={updateVertexIntegration}
                        >
                            Update
                        </Button>
                    </div>
                </Fragment>
            )}
            {!instance?.service_update?.available && (
                <div className={styles.content}>Everything is up-to-date</div>
            )}
        </Vertical>
    );
}
