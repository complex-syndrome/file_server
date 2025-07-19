import type { SettingType } from "./interfaces";

export const SettingDescriptions: SettingType[] = [
    {
        title: "Allow Direct Access of API From Other IPs",
        description: "Allow other IPs to directly connect to the API without having to go through the web UI. (Not including settings)",
        jsonKey: "AllowOtherIPs",
    },
];