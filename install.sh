#!/bin/bash

GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🐱 Installing cutemotd...${NC}"

detect_package_manager() {
    if command -v apt &> /dev/null; then
        echo "apt"
    elif command -v dnf &> /dev/null; then
        echo "dnf"
    elif command -v yum &> /dev/null; then
        echo "yum"
    else
        echo "unknown"
    fi
}

PKG_MANAGER=$(detect_package_manager)

install_dependencies() {
    case $PKG_MANAGER in
        "apt")
            sudo apt update
            sudo apt install -y golang-go git
            ;;
        "dnf")
            sudo dnf install -y golang git
            ;;
        "yum")
            sudo yum install -y golang git
            ;;
        *)
            echo "Unsupported package manager. Please install go and git manually."
            exit 1
            ;;
    esac
}

if ! command -v git &> /dev/null || ! command -v go &> /dev/null; then
    echo -e "${BLUE}📦 Installing dependencies...${NC}"
    install_dependencies
fi

INSTALL_DIR="$HOME/.local/share/cutemotd"
mkdir -p "$INSTALL_DIR"

echo -e "${BLUE}📦 Cloning repository...${NC}"
git clone https://github.com/aubreyrs/cutemotd.git "$INSTALL_DIR/source" || exit 1

cd "$INSTALL_DIR/source"
go build -o cutemotd || exit 1

mkdir -p "$HOME/.local/bin"
mv cutemotd "$HOME/.local/bin/"
chmod +x "$HOME/.local/bin/cutemotd"

if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.bashrc"
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.zshrc" 2>/dev/null || true
fi

MOTD_SCRIPT="$HOME/.local/share/cutemotd/motd.sh"

cat > "$MOTD_SCRIPT" << 'EOF'
#!/bin/bash
if command -v cutemotd &> /dev/null; then
    cutemotd
fi
EOF

chmod +x "$MOTD_SCRIPT"

for rc in "$HOME/.bashrc" "$HOME/.zshrc"; do
    if [ -f "$rc" ]; then
        sed -i '/cutemotd/d' "$rc"
        echo "source $MOTD_SCRIPT" >> "$rc"
    fi
done

echo -e "${GREEN}✨ cutemotd has been installed!${NC}"
echo -e "${BLUE}Please restart your shell or run 'source ~/.bashrc' to see the changes.${NC}"
echo -e "${BLUE}You can test it by running 'cutemotd test'${NC}"
