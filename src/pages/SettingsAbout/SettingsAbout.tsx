import { Fragment } from "react";
import { Title } from "../../components/Text/Text";
import { getAbout } from "../../backend/backend";
import Loading from "../../components/Loading/Loading";
import Symbol from "../../components/Symbol/Symbol";

import styles from "./SettingsAbout.module.sass";
import { Horizontal, Vertical } from "../../components/Layouts/Layouts";
import { useFetch } from "../../hooks/useFetch";
import { About } from "../../models/about";

export default function SettingsAbout() {
    const { data: about, loading, error } = useFetch<About>(getAbout);

    return (
        <Fragment>
            <Title>About</Title>
            {/*{error && <Error error={error} />}*/}
            {loading && !error && <Loading />}
            {!loading && (
                <Vertical gap={4}>
                    <Horizontal gap={12} alignItems="center">
                        <div className={styles.item}>
                            <Symbol name="tag" />
                        </div>
                        <div className={styles.item}>Version:</div>
                        <div className={styles.item}>{about?.version}</div>
                    </Horizontal>
                    <Horizontal gap={12} alignItems="center">
                        <div className={styles.item}>
                            <Symbol name="commit" />
                        </div>
                        <div className={styles.item}>Commit:</div>
                        <div className={styles.item}>{about?.commit}</div>
                    </Horizontal>
                    <Horizontal gap={12} alignItems="center">
                        <div className={styles.item}>
                            <Symbol name="calendar_month" />
                        </div>
                        <div className={styles.item}>Release date:</div>
                        <div className={styles.item}>{about?.date}</div>
                    </Horizontal>
                </Vertical>
            )}
        </Fragment>
    );
}
