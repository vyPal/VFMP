import React, { useState } from 'react';
import { EventsOn } from '../wailsjs/runtime/runtime'
import './App.css';
import { SendCount, SendIndex, SendSearch } from '../wailsjs/go/main/App';

function App() {
    const [selectedOption, setSelectedOption] = useState('');
    const [directory, setDirectory] = useState('');
    const [searchString, setSearchString] = useState('');
    const [fuzzySearch, setFuzzySearch] = useState(false);
    const [fileCount, setFileCount] = useState(0);
    const [indexedFiles, setIndexedFiles] = useState(0);

    const handleSelectChange = (event) => {
        setSelectedOption(event.target.value);
    };

    const handleDirectoryChange = (event) => {
        setDirectory(event.target.value);
    };

    const handleSearchStringChange = (event) => {
        setSearchString(event.target.value);
    };

    const handleFuzzySearchChange = (event) => {
        setFuzzySearch(event.target.checked);
    };

    const handleSubmit = (event) => {
        event.preventDefault();
        if (selectedOption === 'count') {
            SendCount(directory)
        } else if (selectedOption === 'index') {
            SendIndex(directory)
        } else if (selectedOption === 'search') {
            SendSearch("", searchString, fuzzySearch)
        }
    };

    EventsOn('count.progress', (count) => {
        setFileCount(count)
    })

    EventsOn('index.progress', (count) => {
        setIndexedFiles(count)
    })

    EventsOn('search.results', (results) => {
        let resultsDiv = document.getElementById("results")
        resultsDiv.innerHTML = ""
        results.forEach((result) => {
            let resultDiv = document.createElement("div")
            resultDiv.innerHTML = result
            resultsDiv.appendChild(resultDiv)
        })
    })

    EventsOn('seatch.results.fuzzy', (results) => {
        console.log(results)
        let resultsDiv = document.getElementById("results")
        resultsDiv.innerHTML = ""
        results.forEach((result) => {
            console.log()
            let resultDiv = document.createElement("div")
            let highlightedPath = [...result.Path].map((char, index) => 
                result.Indexes.includes(index) ? `<mark>${char}</mark>` : char
            ).join('');
            resultDiv.innerHTML = highlightedPath;
            resultsDiv.appendChild(resultDiv)
        })
    })

    return (
        <div className="App">
            <select value={selectedOption} onChange={handleSelectChange}>
                <option value="">--Please choose an option--</option>
                <option value="count">Count</option>
                <option value="index">Index</option>
                <option value="search">Search</option>
            </select>

            <br />
            <br />

            {(selectedOption === 'count' || selectedOption === 'index') && (
                <input type="text" placeholder="Directory" value={directory} onChange={handleDirectoryChange} />
            )}

            {selectedOption === 'search' && (
                <div>
                    <input type="text" placeholder="Search String" value={searchString} onChange={handleSearchStringChange} />
                    <label>
                        <input type="checkbox" checked={fuzzySearch} onChange={handleFuzzySearchChange} />
                        Fuzzy Search
                    </label>
                </div>
            )}

            {selectedOption && (
                <div>
                    <button onClick={handleSubmit}>Submit</button>
                </div>
            )}

            <br />
            <br />
            <h2>{indexedFiles}/{fileCount}</h2>
            <progress value={indexedFiles} max={fileCount}></progress>
            <div>
                <h2>Results</h2>
                <div id="results"></div>
            </div>
        </div>
    );
}

export default App;