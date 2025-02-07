Raft Protocol Implementation with Hyperledger Fabric
📋 Project Overview
A distributed blockchain network implementing the Raft consensus protocol using Hyperledger Fabric, demonstrating advanced distributed ledger technology with robust fault tolerance.
🚀 Key Technologies

Blockchain: Hyperledger Fabric
Consensus: Raft Protocol
Chaincode: GoLang
UI: Python Flask
Monitoring: Hyperledger Explorer

🔧 System Requirements

Docker 20.10+
Docker Compose
GoLang 1.16+
Python 3.8+
Hyperledger Fabric 2.5

💻 Installation Guide
Prerequisites Setup

Clone Repository
bashCopygit clone <repository-url>
cd Raft_Protocol_Implementation

Generate Crypto Materials
bashCopycd blockchain-network
./generate-artifacts.sh

Start Fabric Network
bashCopydocker-compose up -d

Deploy Chaincode
bashCopy./scripts/deployChaincode.sh

Launch Flask UI
bashCopycd flask-ui
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python app.py

Start Hyperledger Explorer
bashCopycd hyperledger-explorer
./start-explorer.sh


🌐 Network Architecture

Orderer Nodes: 3 Raft consensus nodes
Peer Nodes: 2 organizations with 2 peers each
Channel: Single communication channel
Chaincode: Custom Raft protocol implementation

🔗 Access Points

Blockchain Explorer: http://localhost:8080
Flask UI: http://localhost:5000

🔒 Security Features

TLS encryption
Mutual TLS authentication
Certificate-based access control
Restricted network communication

🛠 Troubleshooting

Verify Docker Containers
bashCopydocker ps

Check Network Logs
bashCopydocker-compose logs


📜 License
Apache 2.0
👥 Contributors

[Your Name]
[Contributor Names]

📞 Contact
Project Maintainer: [Your Name]
Email: [Contact Email]
