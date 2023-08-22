import { Fragment } from "react";
import { Text, Title } from "../../components/Text/Text";
import { useParams } from "react-router-dom";
import useInstance from "../../hooks/useInstance";
import styles from "./BayDetailsUpdate.module.sass";
import { Vertical } from "../../components/Layouts/Layouts";
import Loading from "../../components/Loading/Loading";

export default function BayDetailsUpdate() {
    const { uuid } = useParams();

    const { instance } = useInstance(uuid);

    let content: any;
    if (instance?.update) {
        content = (
            <Vertical gap={10}>
                <Text>An update is available.</Text>
                <Text>Current: {instance?.update?.current_version}</Text>
                <Text>Latest: {instance?.update?.latest_version}</Text>
            </Vertical>
        );
    } else {
        content = (
            <Text className={styles.content}>Everything is up-to-date</Text>
        );
    }

    return (
        <Fragment>
            <Title className={styles.title}>Update</Title>
            {instance ? content : <Loading />}
        </Fragment>
    );
}
