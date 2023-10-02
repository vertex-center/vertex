import styles from "./LoadingLogo.module.sass";
import { Vertical } from "../Layouts/Layouts";
import { SubTitle } from "../Text/Text";

export function FullLoadingLogo({ show }: { show?: boolean }) {
    if (!show) return null;

    return (
        <div className={styles.full}>
            <Vertical gap={10} alignItems="center">
                <LoadingLogo />
                <SubTitle>Loading...</SubTitle>
            </Vertical>
        </div>
    );
}

export default function LoadingLogo() {
    return (
        <div className={styles.logo}>
            <img className={styles.out} src="/images/logo_out.svg" />
            <img className={styles.in} src="/images/logo_in.svg" />
        </div>
    );
}
