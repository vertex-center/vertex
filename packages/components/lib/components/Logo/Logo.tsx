import Icon from "../../assets/logo.svg";
import { ImgHTMLAttributes } from "react";

export type LogoProps = ImgHTMLAttributes<HTMLImageElement> & {
    size?: number;
};

export function Logo(props: Readonly<LogoProps>) {
    const { alt = "App Logo", size = 30, ...others } = props;
    return <img width={size} src={Icon} alt={alt} {...others} />;
}
