package multiaddr

import(
  "fmt"
  "strings"
  "encoding/binary"
  "net"
  "strconv"
)


func StringToBytes(s string) ([]byte, error) {
  b := []byte{}
  sp := strings.Split(s, "/")

  // consume first empty elem
  sp = sp[1:]

  for ; len(sp) > 0 ; {
    p := ProtocolWithName(sp[0])
    if p == nil {
      return nil, fmt.Errorf("no protocol with name %s", sp[0])
    }
    b = append(b, byte(p.Code))

    a := AddressStringToBytes(p, sp[1])
    b = append(b, a...)

    sp = sp[2:]
  }
  return b, nil
}

func BytesToString(b []byte) (string, error) {
  s := ""

  for ; len(b) > 0 ; {
    p := ProtocolWithCode(int(b[0]))
    if p == nil {
      return "", fmt.Errorf("no protocol with code %d", b[0])
    }
    s = strings.Join([]string{s, "/", p.Name}, "")
    b = b[1:]

    a := AddressBytesToString(p, b)
    if len(a) > 0 {
      s = strings.Join([]string{s, "/", a}, "")
    }
    b = b[(p.Size / 8):]
  }

  return s, nil
}

func AddressStringToBytes(p *Protocol, s string) []byte {
  switch p.Code {

  // ipv4,6
  case 4, 41:
    return net.ParseIP(s).To4()

  // tcp udp dccp sctp
  case 6, 17, 33, 132:
    b := make([]byte, 2)
    i, err := strconv.Atoi(s)
    if err == nil {
      binary.BigEndian.PutUint16(b, uint16(i))
    }
    return b
  }

  return []byte{}
}

func AddressBytesToString(p *Protocol, b []byte) string {
  switch p.Code {

  // ipv4,6
  case 4, 41:
    return net.IP(b).String()

  // tcp udp dccp sctp
  case 6, 17, 33, 132:
    i := binary.BigEndian.Uint16(b)
    return strconv.Itoa(int(i))
  }

  return ""
}

