import Progress from "../Progress";
import styles from "./Downloads.module.sass";

export type Download = {
    id: string;
    status: string;
    current?: number;
    total?: number;
};

export default function Downloads(props: { downloads?: Download[] }) {
    const { downloads } = props;

    console.log(downloads);
    if (!downloads) return null;
    if (downloads.length === 0) return null;

    const status = downloads?.find((dl) => dl.id === "")?.status;

    return (
        <div className={styles.downloads}>
            <div className={styles.downloadsList}>
                {downloads.map((dl) => {
                    if (dl.id === "") return null;
                    return (
                        <div key={dl.id} className={styles.download}>
                            {dl.id && (
                                <div className={styles.downloadID}>{dl.id}</div>
                            )}
                            {dl.status}
                            {dl.current && dl.total && (
                                <div className={styles.downloadProgress}>
                                    <Progress
                                        value={(dl.current / dl.total) * 100}
                                    />
                                </div>
                            )}
                        </div>
                    );
                })}
            </div>
            {status && <div className={styles.download}>{status}</div>}
        </div>
    );
}
