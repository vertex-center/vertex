import styles from "./LoadingLogo.module.sass";

export default function LoadingLogo() {
    return (
        <div className={styles.logo}>
            <img className={styles.out} src="/images/logo_out.svg" />
            <img className={styles.in} src="/images/logo_in.svg" />
        </div>
    );
}
