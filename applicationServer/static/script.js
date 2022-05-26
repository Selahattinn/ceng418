const wrapper = document.querySelector(".wrapper"),
searchInput = wrapper.querySelector("input"),
updateWords = wrapper.querySelector(".refreshWordsButton"),
volume = wrapper.querySelector(".word i"),
infoText = wrapper.querySelector(".info-text"),
synonyms = wrapper.querySelector(".synonyms .list"),
removeIcon = wrapper.querySelector(".search span");
let audio;
function data(result, word){
    console.log(result);
    if(result.title){
        infoText.innerHTML = `Can't find the meaning of <span>"${word}"</span>. Please, try to search for another word.`;
    }else{
        wrapper.classList.add("active");
        let definitions = result.meaning
        document.querySelector(".word p").innerText = result.word;
        document.querySelector(".meaning span").innerText = definitions
        if(definitions.synonyms[0] == undefined){
            synonyms.parentElement.style.display = "none";
        }
    }
}
function search(word){
    fetchApi(word);
    searchInput.value = word;
}
function fetchApi(word){
    wrapper.classList.remove("active");
    infoText.style.color = "#000";
    infoText.innerHTML = `Searching the meaning of <span>"${word}"</span>`;
    let url = `http://172.25.152.115:8080/words/${word}`;
    fetch(url).then(response => response.json()).then(result => data(result, word)).catch(() =>{
        infoText.innerHTML = `Can't find the meaning of <span>"${word}"</span>. Please, try to search for another word.`;
    });
}

function updateWordsFunc(){
    let url = `http://172.25.152.115:8080/updateWords`;
    // make post request if it success then refresh
    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    }).then(response => console.log(response)).catch( error => console.log(error));

    
}
searchInput.addEventListener("keyup", e =>{
    let word = e.target.value.replace(/\s+/g, ' ');
    if(e.key == "Enter" && word){
        fetchApi(word);
    }
});
volume.addEventListener("click", ()=>{
    volume.style.color = "#4D59FB";
    audio.play();
    setTimeout(() =>{
        volume.style.color = "#999";
    }, 800);
});
removeIcon.addEventListener("click", ()=>{
    searchInput.value = "";
    searchInput.focus();
    wrapper.classList.remove("active");
    infoText.style.color = "#9A9A9A";
    infoText.innerHTML = "Type any existing word and press enter to get meaning, example, synonyms, etc.";
});

updateWords.addEventListener("click", ()=>{
    updateWordsFunc();
});
