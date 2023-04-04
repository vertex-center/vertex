import { Fragment, useEffect, useState } from "react";
import { Title } from "../../components/Text/Text";
import { getAbout } from "../../backend/backend";
import Loading from "../../components/Loading/Loading";
import Symbol from "../../components/Symbol/Symbol";

import styles from "./SettingsAbout.module.sass";

type Props = {};

export default function SettingsAbout(props: Props) {
    const [version, setVersion] = useState<string>();
    const [commit, setCommit] = useState<string>();
    const [date, setDate] = useState<string>();

    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        setLoading(true);
        getAbout().then((about) => {
            setVersion(about.version);
            setCommit(about.commit);
            setDate(about.date);
            setLoading(false);
        });
    }, []);

    return (
        <Fragment>
            <Title>About</Title>
            {loading && <Loading />}
            {!loading && (
                <table>
                    <tbody>
                        <tr>
                            <td className={styles.item}>
                                <Symbol name="tag" />
                            </td>
                            <td className={styles.item}>Version:</td>
                            <td className={styles.item}>{version}</td>
                        </tr>
                        <tr>
                            <td className={styles.item}>
                                <Symbol name="commit" />
                            </td>
                            <td className={styles.item}>Commit:</td>
                            <td className={styles.item}>{commit}</td>
                        </tr>
                        <tr>
                            <td className={styles.item}>
                                <Symbol name="calendar_month" />
                            </td>
                            <td className={styles.item}>Release date:</td>
                            <td className={styles.item}>{date}</td>
                        </tr>
                    </tbody>
                </table>
            )}
        </Fragment>
    );
}
