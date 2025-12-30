export namespace cache {
	
	export class CacheInfo {
	    downloadCacheSize: number;
	    downloadCachePath: string;
	    totalSize: number;
	    totalSizeHuman: string;
	
	    static createFrom(source: any = {}) {
	        return new CacheInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.downloadCacheSize = source["downloadCacheSize"];
	        this.downloadCachePath = source["downloadCachePath"];
	        this.totalSize = source["totalSize"];
	        this.totalSizeHuman = source["totalSizeHuman"];
	    }
	}

}

export namespace config {
	
	export class Config {
	    mirror: string;
	    proxy: string;
	    goproxy: string;
	    default_version: string;
	    aliases: Record<string, string>;
	    theme: string;
	    language: string;
	    install_dir: string;
	    cache_dir: string;
	    gopath_mode: string;
	    shared_gopath: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mirror = source["mirror"];
	        this.proxy = source["proxy"];
	        this.goproxy = source["goproxy"];
	        this.default_version = source["default_version"];
	        this.aliases = source["aliases"];
	        this.theme = source["theme"];
	        this.language = source["language"];
	        this.install_dir = source["install_dir"];
	        this.cache_dir = source["cache_dir"];
	        this.gopath_mode = source["gopath_mode"];
	        this.shared_gopath = source["shared_gopath"];
	    }
	}
	export class MirrorOption {
	    name: string;
	    url: string;
	
	    static createFrom(source: any = {}) {
	        return new MirrorOption(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.url = source["url"];
	    }
	}

}

export namespace env {
	
	export class DiagnoseResult {
	    item: string;
	    status: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new DiagnoseResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.item = source["item"];
	        this.status = source["status"];
	        this.message = source["message"];
	    }
	}
	export class EnvInfo {
	    goroot: string;
	    gopath: string;
	    goproxy: string;
	    gobin: string;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new EnvInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.goroot = source["goroot"];
	        this.gopath = source["gopath"];
	        this.goproxy = source["goproxy"];
	        this.gobin = source["gobin"];
	        this.path = source["path"];
	    }
	}

}

export namespace main {
	
	export class InstallSettings {
	    installDir: string;
	    cacheDir: string;
	
	    static createFrom(source: any = {}) {
	        return new InstallSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.installDir = source["installDir"];
	        this.cacheDir = source["cacheDir"];
	    }
	}

}

export namespace module {
	
	export class ModuleCacheStats {
	    totalModules: number;
	    totalSize: number;
	    totalSizeStr: string;
	    cachePath: string;
	
	    static createFrom(source: any = {}) {
	        return new ModuleCacheStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalModules = source["totalModules"];
	        this.totalSize = source["totalSize"];
	        this.totalSizeStr = source["totalSizeStr"];
	        this.cachePath = source["cachePath"];
	    }
	}
	export class ModuleInfo {
	    path: string;
	    version: string;
	    size: number;
	    dir: string;
	    name: string;
	    description: string;
	    category: string;
	
	    static createFrom(source: any = {}) {
	        return new ModuleInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.version = source["version"];
	        this.size = source["size"];
	        this.dir = source["dir"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.category = source["category"];
	    }
	}
	export class SearchResult {
	    path: string;
	    version: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new SearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.version = source["version"];
	        this.description = source["description"];
	    }
	}
	export class VerifyResult {
	    module: string;
	    version: string;
	    status: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new VerifyResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.module = source["module"];
	        this.version = source["version"];
	        this.status = source["status"];
	        this.message = source["message"];
	    }
	}

}

export namespace version {
	
	export class FileInfo {
	    filename: string;
	    os: string;
	    arch: string;
	    size: number;
	    sha256: string;
	    kind: string;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filename = source["filename"];
	        this.os = source["os"];
	        this.arch = source["arch"];
	        this.size = source["size"];
	        this.sha256 = source["sha256"];
	        this.kind = source["kind"];
	    }
	}
	export class GoVersion {
	    version: string;
	    stable: boolean;
	    files: FileInfo[];
	
	    static createFrom(source: any = {}) {
	        return new GoVersion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.stable = source["stable"];
	        this.files = this.convertValues(source["files"], FileInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class InstalledVersion {
	    version: string;
	    path: string;
	    isCurrent: boolean;
	
	    static createFrom(source: any = {}) {
	        return new InstalledVersion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.path = source["path"];
	        this.isCurrent = source["isCurrent"];
	    }
	}

}

