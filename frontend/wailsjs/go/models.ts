export namespace userData {
	
	export class HostConfig {
	    id: number;
	    name: string;
	    applyHosts: boolean;
	
	    static createFrom(source: any = {}) {
	        return new HostConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.applyHosts = source["applyHosts"];
	    }
	}

}

