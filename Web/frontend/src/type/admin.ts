export type tempLog = {
    ID: number,
    CreatedAt: string,
    UpdatedAt: string | null,
    DeletedAt: string | null,
    fingerprintData: string,
    userName: string,
    workNum: string
}

export type userInfo = {
    ID: number,
    CreatedAt: string,
    UpdatedAt: string | null,
    DeletedAt: string | null,
    userName: string,
    workNum: string,
    finger1: string,
    finger2: string
}

export type fingerLog = {
    ID: number,
    CreatedAt: string,
    UpdatedAt: string | null,
    DeletedAt: string | null,
    fingerId: string,
    workNum: string,
    status: string
}

export type operateLog = {
    ID: number,
    CreatedAt: string,
    UpdatedAt: string | number,
    DeletedAt: string | number,
    role: string,
    operate: string,
    object: string
}