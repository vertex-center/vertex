import "./ProfilePicture.sass";
import cx from "classnames";
import { HTMLProps } from "react";

export type ProfilePictureProps = HTMLProps<HTMLImageElement> & {
    size?: number;
};

export function ProfilePicture(props: Readonly<ProfilePictureProps>) {
    const { className, alt, size = 40, ...others } = props;

    return (
        <img
            className={cx("profile-picture", className)}
            alt={alt}
            width={size}
            height={size}
            {...others}
        />
    );
}
