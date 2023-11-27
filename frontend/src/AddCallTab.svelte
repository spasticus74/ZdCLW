<script>
    //import { SearchContactByName, SearchContactByOrg } from "../wailsjs/go/main/App.js";
    //import Select from "./Select.svelte";
    import SiteSelect from "./OrgSelector.svelte";
    import ProductSelect from "./ProductSelector.svelte";

    let contacts = [];
    let resultsVisible = false

    let callerName = "";
    let orgName = "";
    let subject = ""; //issue
    let problem = "";
    let assetName = "";
    let product = "";
    let phoneNumber = "";
    let issueStartDate = "";

    function submitNewTicket() {
        resultsVisible = false
        contacts =  [];
    }

</script>
<h2>
    Create a Zendesk ticket
</h2>

<table align="center">
    <tr>
        <td align="right"><label for="org">Organisation</label></td>
        <td align="left"><SiteSelect bind:value={orgName}/></td>
    </tr>
    <tr>
        <td align="right"><label for="caller">Caller</label></td>
        <td align="left"><input id="caller" bind:value={callerName}></td>
    </tr>
    <tr>
        <td align="right"><label for="subj">Subject</label></td>
        <td align="left"><input id="subj" bind:value={subject}></td>
    </tr>
    <tr>
        <td align="right"><label for="prob">Problem</label></td>
        <td align="left"><textarea id="prob" bind:value={problem} /></td>
    </tr>
    <tr>
        <td align="right"><label for="asset">Asset</label></td>
        <td align="left"><input id="asset" bind:value={assetName} /></td>
    </tr>
    <tr>
        <td align="right"><label for="prod">Product</label></td>
        <td align="left"><ProductSelect  display_func={o => o.text} bind:value={product}/></td>
    </tr>
    <tr>
        <td align="right"><label for="phone">Caller's phone #</label></td>
        <td align="left"><input id="phone" bind:value={phoneNumber} /></td>
    </tr>
    <tr>
        <td align="right"><label for="startDate">Date issue began</label></td>
        <td align="left"><input id="startDate" bind:value={issueStartDate} /></td>
    </tr>
</table>
    
<button on:click={submitNewTicket}>
    Submit
</button>   

    
<style>
    button {
        border-radius: .25rem;
        cursor: pointer;
        font-weight: bold;
        background-color: #fff;
        box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
        max-width: 250px;
    }
    button:active {
        background-color: #efefef;
        box-shadow: 0 0px 0px 0 rgba(0, 0, 0, 0.1), 0 0px 0px 0 rgba(0, 0, 0, 0.06);	
    }
    textarea {
		width: 100%;
		height: 150px;
	}
</style>