export namespace service {
	
	export class HostConfig {
	    id: number;
	    name: string;
	    applyHosts: boolean;
	    defaultTarget: string;
	    enableTLS: boolean;
	    TLSCertFile: string;
	    TLSKeyFile: string;
	
	    static createFrom(source: any = {}) {
	        return new HostConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.applyHosts = source["applyHosts"];
	        this.defaultTarget = source["defaultTarget"];
	        this.enableTLS = source["enableTLS"];
	        this.TLSCertFile = source["TLSCertFile"];
	        this.TLSKeyFile = source["TLSKeyFile"];
	    }
	}

}

