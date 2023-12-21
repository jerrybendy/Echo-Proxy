export namespace service {
	
	export class HostProxy {
	    id: number;
	    matchType: string;
	    matchRule: string;
	    matchParams: {[key: string]: any};
	    targetType: string;
	    targetParams: {[key: string]: any};
	
	    static createFrom(source: any = {}) {
	        return new HostProxy(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.matchType = source["matchType"];
	        this.matchRule = source["matchRule"];
	        this.matchParams = source["matchParams"];
	        this.targetType = source["targetType"];
	        this.targetParams = source["targetParams"];
	    }
	}
	export class HostConfig {
	    id: number;
	    name: string;
	    applyHosts: boolean;
	    enableTLS: boolean;
	    TLSCertFile: string;
	    TLSKeyFile: string;
	    proxies: HostProxy[];
	
	    static createFrom(source: any = {}) {
	        return new HostConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.applyHosts = source["applyHosts"];
	        this.enableTLS = source["enableTLS"];
	        this.TLSCertFile = source["TLSCertFile"];
	        this.TLSKeyFile = source["TLSKeyFile"];
	        this.proxies = this.convertValues(source["proxies"], HostProxy);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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
	export class Setting {
	    httpPort: number;
	    httpsPort: number;
	
	    static createFrom(source: any = {}) {
	        return new Setting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.httpPort = source["httpPort"];
	        this.httpsPort = source["httpsPort"];
	    }
	}

}

