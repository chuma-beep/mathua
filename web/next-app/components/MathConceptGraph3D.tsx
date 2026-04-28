'use client'

import React, { useRef, useState, useCallback, useMemo, useEffect } from 'react'
import { Canvas, useFrame, useThree } from '@react-three/fiber'
import { OrbitControls, Html } from '@react-three/drei'
import * as THREE from 'three'

export type MathField =
  | "Number Theory" | "Algebra" | "Geometry" | "Linear Algebra"
  | "Discrete Math" | "Calculus" | "Analysis" | "Topology" | "Abstract Algebra"

export interface MathNode {
  id: string
  name: string
  field: MathField
  description: string
  level: "beginner" | "intermediate" | "advanced"
  position: [number, number, number]
}

export interface MathLink {
  source: string
  target: string
  type: "prerequisite" | "related" | "extends"
}

interface MathConceptGraph3DProps {
  theme?: 'dark' | 'light'
}

const FIELD_COLORS: Record<MathField, string> = {
  "Number Theory":    "#c8a96e",
  "Algebra":          "#4db8a0",
  "Geometry":         "#a8a0f0",
  "Linear Algebra":   "#7dd3fc",
  "Discrete Math":    "#86efac",
  "Calculus":         "#fda4af",
  "Analysis":         "#f0abfc",
  "Topology":         "#fb923c",
  "Abstract Algebra": "#e879f9",
}

const FIELD_COLORS_LIGHT: Record<MathField, string> = {
  "Number Theory":    "#a0814a",
  "Algebra":          "#3a9a8a",
  "Geometry":         "#8a80d8",
  "Linear Algebra":   "#5ab8dc",
  "Discrete Math":    "#68c88c",
  "Calculus":         "#e88a99",
  "Analysis":         "#c88ae8",
  "Topology":         "#e87830",
  "Abstract Algebra": "#c868e8",
}

const LINK_COLORS: Record<MathLink["type"], string> = {
  prerequisite: "#c8a96e",
  related:      "#4db8a0",
  extends:      "#a8a0f0",
}

const LINK_COLORS_LIGHT: Record<MathLink["type"], string> = {
  prerequisite: "#a0814a",
  related:      "#3a9a8a",
  extends:      "#8a80d8",
}

const NODE_RADIUS: Record<MathNode["level"], number> = {
  beginner:     0.22,
  intermediate: 0.18,
  advanced:     0.15,
}

const NODES: MathNode[] = [
  { id: "counting",            name: "Counting",               field: "Number Theory",   level: "beginner",     position: [-14.0,  8.0, -4.0], description: "Cardinality, ordering, and the number line from zero." },
  { id: "arithmetic",          name: "Arithmetic",             field: "Number Theory",   level: "beginner",     position: [-11.0,  6.0, -1.0], description: "Addition, subtraction, multiplication, and division of integers." },
  { id: "number-theory",       name: "Number Theory",          field: "Number Theory",   level: "beginner",     position: [-11.0,  7.2, -2.5], description: "Primes, divisibility, modular arithmetic, and integer structure." },
  { id: "fractions",           name: "Fractions",              field: "Number Theory",   level: "beginner",     position: [-13.0,  4.0,  2.0], description: "Rational numbers, equivalence, and operations on parts of a whole." },
  { id: "decimals",            name: "Decimals",               field: "Number Theory",   level: "beginner",     position: [-12.0,  2.0,  3.0], description: "Base-10 representation of non-integer quantities." },
  { id: "percentages",         name: "Percentages",            field: "Number Theory",   level: "beginner",     position: [-10.0,  1.0,  5.0], description: "Ratios expressed as parts per hundred." },
  { id: "ratios",              name: "Ratios and Proportions", field: "Number Theory",   level: "beginner",     position: [-11.5,  3.0,  1.0], description: "Multiplicative relationships and scaling between quantities." },
  { id: "negative-numbers",    name: "Negative Numbers",       field: "Number Theory",   level: "beginner",     position: [-10.5,  5.5,  1.5], description: "Extension of the number line below zero, signed arithmetic." },
  { id: "exponents",           name: "Exponents and Roots",    field: "Number Theory",   level: "beginner",     position: [ -9.0,  4.5,  3.5], description: "Repeated multiplication, powers, and inverse root operations." },
  { id: "order-of-operations", name: "Order of Operations",    field: "Number Theory",   level: "beginner",     position: [ -9.5,  6.5, -2.0], description: "BODMAS rules governing evaluation order in expressions." },
  { id: "variables",           name: "Variables",              field: "Algebra",         level: "beginner",     position: [ -7.0,  4.0,  5.0], description: "Symbolic placeholders representing unknown or changing quantities." },
  { id: "expressions",         name: "Algebraic Expressions",  field: "Algebra",         level: "beginner",     position: [ -5.5,  2.8,  6.0], description: "Combining variables, constants, and operations into symbolic forms." },
  { id: "linear-equations",    name: "Linear Equations",       field: "Algebra",         level: "beginner",     position: [ -4.0,  1.8,  6.8], description: "Degree-one equations — solving for an unknown in one variable." },
  { id: "inequalities",        name: "Inequalities",           field: "Algebra",         level: "beginner",     position: [ -3.0,  3.0,  7.5], description: "Order relations between expressions and solution sets." },
  { id: "systems-of-equations",name: "Systems of Equations",   field: "Algebra",         level: "intermediate", position: [ -2.0,  1.0,  8.0], description: "Simultaneous equations solved by substitution or elimination." },
  { id: "algebra",             name: "Algebra",                field: "Algebra",         level: "beginner",     position: [ -4.4,  1.4,  6.4], description: "Symbolic manipulation, expressions, equations, and functions." },
  { id: "polynomials",         name: "Polynomials",            field: "Algebra",         level: "intermediate", position: [ -1.5,  0.0,  7.2], description: "Sums of power terms — addition, multiplication, factoring, and roots." },
  { id: "quadratics",          name: "Quadratic Equations",    field: "Algebra",         level: "intermediate", position: [  0.5, -0.8,  7.8], description: "Degree-two equations solved by factoring or the quadratic formula." },
  { id: "functions",           name: "Functions",              field: "Algebra",         level: "intermediate", position: [  2.0,  1.0,  7.0], description: "Mappings from inputs to outputs — domain, range, and inverses." },
  { id: "exponential-functions",name:"Exponential Functions",  field: "Algebra",         level: "intermediate", position: [  3.5,  2.0,  6.2], description: "Functions of the form aˣ modelling growth, decay, compounding." },
  { id: "logarithms",          name: "Logarithms",             field: "Algebra",         level: "intermediate", position: [  5.0,  2.8,  5.5], description: "Inverse of exponentiation — solving for the exponent." },
  { id: "sequences-series",    name: "Sequences and Series",   field: "Algebra",         level: "intermediate", position: [ -2.5,  2.5,  8.2], description: "Ordered lists and their sums — arithmetic, geometric, convergence." },
  { id: "geometry",            name: "Geometry",               field: "Geometry",        level: "beginner",     position: [  8.0,  5.8, -6.0], description: "Shape, measurement, spatial structure, and proof." },
  { id: "coordinate-geometry", name: "Coordinate Geometry",    field: "Geometry",        level: "intermediate", position: [  6.5,  4.2, -4.0], description: "Algebraic description of geometric objects in the Cartesian plane." },
  { id: "trigonometry",        name: "Trigonometry",           field: "Geometry",        level: "intermediate", position: [ 12.2,  9.2,  4.2], description: "Angles, periodic functions, identities, and triangle relationships." },
  { id: "vectors",             name: "Vectors",                field: "Geometry",        level: "intermediate", position: [  7.0,  3.5,  3.0], description: "Directed quantities with magnitude and direction." },
  { id: "conic-sections",      name: "Conic Sections",         field: "Geometry",        level: "intermediate", position: [  9.0,  5.0, -2.0], description: "Parabolas, ellipses, hyperbolas as plane-cone intersections." },
  { id: "solid-geometry",      name: "Solid Geometry",         field: "Geometry",        level: "intermediate", position: [ 10.0,  7.0, -5.5], description: "Volume, surface area, and properties of 3D figures." },
  { id: "linear-algebra",      name: "Linear Algebra",         field: "Linear Algebra",  level: "intermediate", position: [  5.2,  0.0,  7.2], description: "Vectors, matrices, transformations, and finite-dimensional spaces." },
  { id: "matrices",            name: "Matrices",               field: "Linear Algebra",  level: "intermediate", position: [  6.8, -0.8,  6.5], description: "Rectangular arrays encoding linear transformations and data." },
  { id: "determinants",        name: "Determinants",           field: "Linear Algebra",  level: "intermediate", position: [  8.2, -1.4,  5.5], description: "Scalar encoding of a matrix — invertibility and volume scaling." },
  { id: "eigenvalues",         name: "Eigenvalues",            field: "Linear Algebra",  level: "advanced",     position: [  9.5, -2.0,  4.2], description: "Invariant directions and scaling factors of a linear transformation." },
  { id: "vector-spaces",       name: "Vector Spaces",          field: "Linear Algebra",  level: "advanced",     position: [  4.0, -1.2,  8.0], description: "Abstract spaces with addition and scalar multiplication." },
  { id: "discrete-math",       name: "Discrete Math",          field: "Discrete Math",   level: "intermediate", position: [ -9.8, -5.4,  1.2], description: "Logic, counting, graphs, recurrence, and finite structures." },
  { id: "logic",               name: "Propositional Logic",    field: "Discrete Math",   level: "intermediate", position: [ -8.5, -4.0,  2.5], description: "Truth values, connectives, tautologies, and logical inference." },
  { id: "set-theory",          name: "Set Theory",             field: "Discrete Math",   level: "intermediate", position: [ -7.5, -3.0,  3.5], description: "Collections, unions, intersections, functions, and cardinality." },
  { id: "combinatorics",       name: "Combinatorics",          field: "Discrete Math",   level: "intermediate", position: [-13.2,-10.2, -3.4], description: "Counting arrangements, structures, selections, and finite patterns." },
  { id: "graph-theory",        name: "Graph Theory",           field: "Discrete Math",   level: "intermediate", position: [ -7.4,-13.2, -8.6], description: "Networks, paths, connectivity, trees, and structural relationships." },
  { id: "recurrence",          name: "Recurrence Relations",   field: "Discrete Math",   level: "intermediate", position: [-11.0, -8.0, -1.0], description: "Sequences defined in terms of previous terms." },
  { id: "probability",         name: "Probability",            field: "Discrete Math",   level: "intermediate", position: [-11.5, -6.5,  3.0], description: "Likelihood of events, distributions, and expected values." },
  { id: "calculus",            name: "Calculus",               field: "Calculus",        level: "intermediate", position: [  0.8, -7.6,  3.4], description: "Limits, derivatives, integrals, and the mathematics of change." },
  { id: "limits",              name: "Limits",                 field: "Calculus",        level: "intermediate", position: [ -0.5, -6.0,  4.8], description: "Behaviour of functions as inputs approach a value." },
  { id: "derivatives",         name: "Derivatives",            field: "Calculus",        level: "intermediate", position: [  1.0, -9.0,  5.0], description: "Instantaneous rate of change — rules, chain rule, applications." },
  { id: "integrals",           name: "Integrals",              field: "Calculus",        level: "intermediate", position: [  2.5, -9.5,  3.8], description: "Accumulation — definite, indefinite, and the fundamental theorem." },
  { id: "multivariable-calculus",name:"Multivariable Calculus",field: "Calculus",        level: "advanced",     position: [  4.5,-10.0,  2.0], description: "Partial derivatives, gradients, multiple integrals, vector fields." },
  { id: "differential-equations",name:"Differential Equations",field:"Calculus",         level: "advanced",     position: [ 11.2,-11.4,  5.8], description: "Equations describing change through derivatives and dynamic systems." },
  { id: "analysis",            name: "Real Analysis",          field: "Analysis",        level: "advanced",     position: [  9.2, -6.6,  0.2], description: "Rigorous limits, continuity, convergence, and real functions." },
  { id: "complex-analysis",    name: "Complex Analysis",       field: "Analysis",        level: "advanced",     position: [  4.4,-12.8, -5.6], description: "Differentiation and integration over complex-valued functions." },
  { id: "functional-analysis", name: "Functional Analysis",    field: "Analysis",        level: "advanced",     position: [ 11.0, -8.0, -2.0], description: "Infinite-dimensional vector spaces, operators, spectral theory." },
  { id: "measure-theory",      name: "Measure Theory",         field: "Analysis",        level: "advanced",     position: [  7.5,-11.0, -4.0], description: "Rigorous foundations of integration via sigma-algebras." },
  { id: "fourier-analysis",    name: "Fourier Analysis",       field: "Analysis",        level: "advanced",     position: [  6.0,-12.0,  1.0], description: "Decomposition of functions into frequency components." },
  { id: "topology",            name: "Topology",               field: "Topology",        level: "advanced",     position: [ 11.2,  1.6, -1.2], description: "Continuity, spaces, compactness, invariants under deformation." },
  { id: "metric-spaces",       name: "Metric Spaces",          field: "Topology",        level: "advanced",     position: [ 10.5, -0.5,  1.5], description: "Sets with a distance function — convergence, completeness." },
  { id: "algebraic-topology",  name: "Algebraic Topology",     field: "Topology",        level: "advanced",     position: [ 12.5,  3.0, -3.5], description: "Homotopy, homology, and fundamental groups." },
  { id: "differential-geometry",name:"Differential Geometry",  field: "Topology",        level: "advanced",     position: [ 13.0,  4.5,  2.0], description: "Smooth manifolds, curvature, geodesics, Riemannian metrics." },
  { id: "abstract-algebra",    name: "Abstract Algebra",       field: "Abstract Algebra",level: "advanced",     position: [ -1.8, -1.6, -9.2], description: "Groups, rings, fields, and algebraic structures." },
  { id: "group-theory",        name: "Group Theory",           field: "Abstract Algebra",level: "advanced",     position: [ -3.0, -2.8,-10.0], description: "Sets with a binary operation satisfying closure, identity, inverses." },
  { id: "ring-theory",         name: "Ring Theory",            field: "Abstract Algebra",level: "advanced",     position: [ -1.0, -3.5,-10.8], description: "Algebraic structures with addition and multiplication." },
  { id: "field-theory",        name: "Field Theory",           field: "Abstract Algebra",level: "advanced",     position: [  1.0, -4.2,-10.0], description: "Rings where every non-zero element has a multiplicative inverse." },
  { id: "galois-theory",       name: "Galois Theory",          field: "Abstract Algebra",level: "advanced",     position: [  0.5, -5.0,-11.2], description: "Symmetries of polynomial roots linking group theory and fields." },
  { id: "category-theory",     name: "Category Theory",        field: "Abstract Algebra",level: "advanced",     position: [  3.0, -2.0, -9.5], description: "Objects, morphisms, functors — abstract structure across mathematics." },
]

const LINKS: MathLink[] = [
  { source: "counting",             target: "arithmetic",           type: "prerequisite" },
  { source: "arithmetic",           target: "fractions",            type: "prerequisite" },
  { source: "arithmetic",           target: "negative-numbers",     type: "prerequisite" },
  { source: "arithmetic",           target: "exponents",            type: "prerequisite" },
  { source: "arithmetic",           target: "order-of-operations",  type: "prerequisite" },
  { source: "fractions",            target: "decimals",             type: "prerequisite" },
  { source: "fractions",            target: "ratios",               type: "prerequisite" },
  { source: "decimals",             target: "percentages",          type: "prerequisite" },
  { source: "ratios",               target: "percentages",          type: "related" },
  { source: "arithmetic",           target: "variables",            type: "prerequisite" },
  { source: "order-of-operations",  target: "expressions",          type: "prerequisite" },
  { source: "variables",            target: "expressions",          type: "prerequisite" },
  { source: "expressions",          target: "linear-equations",     type: "prerequisite" },
  { source: "linear-equations",     target: "inequalities",         type: "extends" },
  { source: "linear-equations",     target: "systems-of-equations", type: "extends" },
  { source: "number-theory",        target: "algebra",              type: "prerequisite" },
  { source: "linear-equations",     target: "algebra",              type: "prerequisite" },
  { source: "algebra",              target: "polynomials",          type: "prerequisite" },
  { source: "polynomials",          target: "quadratics",           type: "prerequisite" },
  { source: "algebra",              target: "functions",            type: "prerequisite" },
  { source: "quadratics",           target: "functions",            type: "related" },
  { source: "exponents",            target: "exponential-functions",type: "prerequisite" },
  { source: "functions",            target: "exponential-functions",type: "prerequisite" },
  { source: "exponential-functions",target: "logarithms",           type: "extends" },
  { source: "algebra",              target: "sequences-series",     type: "prerequisite" },
  { source: "functions",            target: "sequences-series",     type: "related" },
  { source: "arithmetic",           target: "geometry",             type: "prerequisite" },
  { source: "algebra",              target: "coordinate-geometry",  type: "prerequisite" },
  { source: "geometry",             target: "coordinate-geometry",  type: "extends" },
  { source: "geometry",             target: "trigonometry",         type: "prerequisite" },
  { source: "coordinate-geometry",  target: "trigonometry",         type: "related" },
  { source: "algebra",              target: "vectors",              type: "prerequisite" },
  { source: "geometry",             target: "vectors",              type: "related" },
  { source: "algebra",              target: "conic-sections",       type: "prerequisite" },
  { source: "coordinate-geometry",  target: "conic-sections",       type: "extends" },
  { source: "geometry",             target: "solid-geometry",       type: "prerequisite" },
  { source: "algebra",              target: "linear-algebra",       type: "prerequisite" },
  { source: "geometry",             target: "linear-algebra",       type: "related" },
  { source: "vectors",              target: "linear-algebra",       type: "prerequisite" },
  { source: "linear-algebra",       target: "matrices",             type: "prerequisite" },
  { source: "matrices",             target: "determinants",         type: "prerequisite" },
  { source: "determinants",         target: "eigenvalues",          type: "prerequisite" },
  { source: "linear-algebra",       target: "vector-spaces",        type: "extends" },
  { source: "abstract-algebra",     target: "vector-spaces",        type: "related" },
  { source: "arithmetic",           target: "discrete-math",        type: "prerequisite" },
  { source: "algebra",              target: "logic",                type: "prerequisite" },
  { source: "discrete-math",        target: "logic",                type: "prerequisite" },
  { source: "logic",                target: "set-theory",           type: "prerequisite" },
  { source: "arithmetic",           target: "combinatorics",        type: "prerequisite" },
  { source: "discrete-math",        target: "combinatorics",        type: "extends" },
  { source: "set-theory",           target: "combinatorics",        type: "related" },
  { source: "combinatorics",        target: "graph-theory",         type: "related" },
  { source: "discrete-math",        target: "recurrence",           type: "extends" },
  { source: "sequences-series",     target: "recurrence",           type: "related" },
  { source: "combinatorics",        target: "probability",          type: "prerequisite" },
  { source: "fractions",            target: "probability",          type: "prerequisite" },
  { source: "number-theory",        target: "combinatorics",        type: "related" },
  { source: "discrete-math",        target: "abstract-algebra",     type: "extends" },
  { source: "number-theory",        target: "abstract-algebra",     type: "extends" },
  { source: "algebra",              target: "calculus",             type: "prerequisite" },
  { source: "trigonometry",         target: "calculus",             type: "prerequisite" },
  { source: "functions",            target: "limits",               type: "prerequisite" },
  { source: "calculus",             target: "limits",               type: "prerequisite" },
  { source: "limits",               target: "derivatives",          type: "prerequisite" },
  { source: "derivatives",          target: "integrals",            type: "prerequisite" },
  { source: "linear-algebra",       target: "multivariable-calculus",type:"prerequisite" },
  { source: "integrals",            target: "multivariable-calculus",type:"prerequisite" },
  { source: "calculus",             target: "differential-equations",type:"extends" },
  { source: "linear-algebra",       target: "differential-equations",type:"related" },
  { source: "multivariable-calculus",target:"differential-equations",type:"related" },
  { source: "sequences-series",     target: "limits",               type: "related" },
  { source: "calculus",             target: "analysis",             type: "extends" },
  { source: "metric-spaces",        target: "analysis",             type: "prerequisite" },
  { source: "analysis",             target: "complex-analysis",     type: "extends" },
  { source: "complex-analysis",     target: "differential-equations",type:"related" },
  { source: "analysis",             target: "functional-analysis",  type: "extends" },
  { source: "vector-spaces",        target: "functional-analysis",  type: "related" },
  { source: "integrals",            target: "measure-theory",       type: "prerequisite" },
  { source: "analysis",             target: "measure-theory",       type: "extends" },
  { source: "measure-theory",       target: "functional-analysis",  type: "related" },
  { source: "integrals",            target: "fourier-analysis",     type: "prerequisite" },
  { source: "analysis",             target: "fourier-analysis",     type: "extends" },
  { source: "complex-analysis",     target: "fourier-analysis",     type: "related" },
  { source: "analysis",             target: "topology",             type: "extends" },
  { source: "geometry",             target: "topology",             type: "related" },
  { source: "graph-theory",         target: "topology",             type: "related" },
  { source: "metric-spaces",        target: "topology",             type: "prerequisite" },
  { source: "analysis",             target: "metric-spaces",        type: "prerequisite" },
  { source: "topology",             target: "algebraic-topology",   type: "extends" },
  { source: "abstract-algebra",     target: "algebraic-topology",   type: "related" },
  { source: "calculus",             target: "differential-geometry",type: "prerequisite" },
  { source: "linear-algebra",       target: "differential-geometry",type: "prerequisite" },
  { source: "topology",             target: "differential-geometry",type: "extends" },
  { source: "linear-algebra",       target: "abstract-algebra",     type: "related" },
  { source: "abstract-algebra",     target: "group-theory",         type: "prerequisite" },
  { source: "group-theory",         target: "ring-theory",          type: "extends" },
  { source: "ring-theory",          target: "field-theory",         type: "extends" },
  { source: "field-theory",         target: "galois-theory",        type: "prerequisite" },
  { source: "group-theory",         target: "galois-theory",        type: "prerequisite" },
  { source: "abstract-algebra",     target: "category-theory",      type: "extends" },
  { source: "topology",             target: "category-theory",      type: "related" },
  { source: "algebraic-topology",   target: "category-theory",      type: "related" },
]

const nodeMap = new Map(NODES.map(n => [n.id, n]))

interface GraphSceneProps {
  nodes: MathNode[]
  links: MathLink[]
  activeId: string
  onSelect: (id: string) => void
  theme: 'dark' | 'light'
}

const NodeMesh = React.memo(function NodeMesh({ 
  node, 
  isActive, 
  isHovered, 
  onHover,
  onSelect,
  theme 
}: { 
  node: MathNode
  isActive: boolean
  isHovered: boolean
  onHover: (id: string | null) => void
  onSelect: () => void
  theme: 'dark' | 'light'
}) {
  const meshRef = useRef<THREE.Mesh>(null)
  const glowRef = useRef<THREE.Mesh>(null)
  const radius = NODE_RADIUS[node.level]
  const fieldColors = theme === 'dark' ? FIELD_COLORS : FIELD_COLORS_LIGHT
  const fieldColor = fieldColors[node.field]
  const [isMobile, setIsMobile] = useState(false)

  useEffect(() => {
    const checkMobile = () => setIsMobile(window.innerWidth < 640)
    checkMobile()
    window.addEventListener('resize', checkMobile)
    return () => window.removeEventListener('resize', checkMobile)
  }, [])
  
  useFrame(() => {
    if (meshRef.current) {
      const mat = meshRef.current.material as THREE.MeshStandardMaterial
      const targetIntensity = isActive ? 2.0 : isHovered ? 1.2 : 0.3
      mat.emissiveIntensity = THREE.MathUtils.lerp(mat.emissiveIntensity, targetIntensity, 0.12)
    }
    if (glowRef.current) {
      const mat = glowRef.current.material as THREE.MeshBasicMaterial
      const targetOpacity = isHovered ? 0.12 : 0.04
      mat.opacity = THREE.MathUtils.lerp(mat.opacity, targetOpacity, 0.1)
    }
  })

  const textColor = theme === 'dark' ? '#c8a96e' : '#a0814a'
  const labelColor = theme === 'dark' ? '#5a6577' : '#888'
  const descColor = theme === 'dark' ? '#9fa8b4' : '#666'

  return (
    <group position={node.position as [number, number, number]}>
      <mesh
        ref={glowRef}
        onPointerOver={() => onHover(node.id)}
        onPointerOut={() => onHover(null)}
        onClick={onSelect}
      >
        <sphereGeometry args={[radius * 2.2, 12, 12]} />
        <meshBasicMaterial color={fieldColor} transparent opacity={0.04} depthWrite={false} />
      </mesh>
      <mesh ref={meshRef}>
        <sphereGeometry args={[radius, 16, 16]} />
        <meshStandardMaterial
          color={fieldColor}
          emissive={fieldColor}
          emissiveIntensity={0.3}
          roughness={0.2}
          metalness={0.1}
        />
      </mesh>
      {isHovered && (
        <Html center distanceFactor={isMobile ? 18 : 12} style={{ pointerEvents: 'none', zIndex: 1000 }}>
          <div style={{
            background: theme === 'dark' ? 'rgba(11,15,26,0.98)' : 'rgba(255,255,255,0.98)',
            border: `0.5px solid ${theme === 'dark' ? '#2a3f5f' : '#ddd'}`,
            borderRadius: '8px',
            padding: isMobile ? '10px 14px' : '8px 12px',
            whiteSpace: 'nowrap',
            minWidth: 'max-content',
            maxWidth: isMobile ? '220px' : '240px',
            boxShadow: theme === 'dark' ? '0 4px 20px rgba(0,0,0,0.5)' : '0 4px 20px rgba(0,0,0,0.15)'
          }}>
            <div style={{ color: textColor, fontFamily: 'monospace', fontSize: isMobile ? '14px' : '12px', fontWeight: 600 }}>{node.name}</div>
            <div style={{ color: labelColor, fontFamily: 'monospace', fontSize: isMobile ? '11px' : '10px', textTransform: 'uppercase', marginTop: '3px' }}>{node.field}</div>
            <div style={{ color: fieldColor, fontFamily: 'monospace', fontSize: isMobile ? '11px' : '10px', textTransform: 'uppercase', marginTop: '2px' }}>{node.level}</div>
            <div style={{ color: descColor, fontFamily: 'monospace', fontSize: isMobile ? '12px' : '11px', maxWidth: isMobile ? '200px' : '200px', whiteSpace: 'normal', marginTop: '6px', lineHeight: '1.4' }}>{node.description}</div>
          </div>
        </Html>
      )}
    </group>
  )
})

function EdgeLines({ links, nodes, activeId, theme }: { links: MathLink[], nodes: MathNode[], activeId: string, theme: 'dark' | 'light' }) {
  const linkColors = theme === 'dark' ? LINK_COLORS : LINK_COLORS_LIGHT

  const activeLinks = useMemo(() => {
    return links.filter(l => l.source === activeId || l.target === activeId)
  }, [links, activeId])

  const geometry = useMemo(() => {
    const positions = new Float32Array(activeLinks.length * 6)
    const colors = new Float32Array(activeLinks.length * 6)
    
    activeLinks.forEach((link, i) => {
      const sourceNode = nodeMap.get(link.source)
      const targetNode = nodeMap.get(link.target)
      if (!sourceNode || !targetNode) return
      
      const idx = i * 6
      positions[idx] = sourceNode.position[0]
      positions[idx+1] = sourceNode.position[1]
      positions[idx+2] = sourceNode.position[2]
      positions[idx+3] = targetNode.position[0]
      positions[idx+4] = targetNode.position[1]
      positions[idx+5] = targetNode.position[2]
      
      const color = new THREE.Color(linkColors[link.type])
      colors[idx] = color.r
      colors[idx+1] = color.g
      colors[idx+2] = color.b
      colors[idx+3] = color.r
      colors[idx+4] = color.g
      colors[idx+5] = color.b
    })
    
    const geo = new THREE.BufferGeometry()
    geo.setAttribute('position', new THREE.BufferAttribute(positions, 3))
    geo.setAttribute('color', new THREE.BufferAttribute(colors, 3))
    return geo
  }, [activeLinks, linkColors])

  const lineRef = useRef<THREE.LineSegments>(null)
  const clockRef = useRef(0)
  
  useFrame((_, delta) => {
    clockRef.current += delta
    if (lineRef.current) {
      const mat = lineRef.current.material as THREE.LineBasicMaterial
      const pulse = 0.6 + 0.4 * Math.sin(clockRef.current * 2)
      mat.opacity = pulse
    }
  })

  if (activeLinks.length === 0) return null

  return (
    <lineSegments ref={lineRef} geometry={geometry}>
      <lineBasicMaterial vertexColors transparent opacity={0.8} linewidth={2} />
    </lineSegments>
  )
}

function AllEdges({ links, theme }: { links: MathLink[], theme: 'dark' | 'light' }) {
  const linkColors = theme === 'dark' ? LINK_COLORS : LINK_COLORS_LIGHT

  const geometry = useMemo(() => {
    const positions = new Float32Array(links.length * 6)
    const colors = new Float32Array(links.length * 6)
    
    links.forEach((link, i) => {
      const sourceNode = nodeMap.get(link.source)
      const targetNode = nodeMap.get(link.target)
      if (!sourceNode || !targetNode) return
      
      const idx = i * 6
      positions[idx] = sourceNode.position[0]
      positions[idx+1] = sourceNode.position[1]
      positions[idx+2] = sourceNode.position[2]
      positions[idx+3] = targetNode.position[0]
      positions[idx+4] = targetNode.position[1]
      positions[idx+5] = targetNode.position[2]
      
      const color = new THREE.Color(linkColors[link.type])
      colors[idx] = color.r
      colors[idx+1] = color.g
      colors[idx+2] = color.b
      colors[idx+3] = color.r
      colors[idx+4] = color.g
      colors[idx+5] = color.b
    })
    
    const geo = new THREE.BufferGeometry()
    geo.setAttribute('position', new THREE.BufferAttribute(positions, 3))
    geo.setAttribute('color', new THREE.BufferAttribute(colors, 3))
    return geo
  }, [links, linkColors])

  return (
    <lineSegments geometry={geometry}>
      <lineBasicMaterial vertexColors transparent opacity={theme === 'dark' ? 0.15 : 0.12} />
    </lineSegments>
  )
}

function Particles({ theme }: { theme: 'dark' | 'light' }) {
  const positions = useMemo(() => {
    const pos = new Float32Array(800 * 3)
    for (let i = 0; i < 800; i++) {
      const theta = Math.random() * Math.PI * 2
      const phi = Math.acos(2 * Math.random() - 1)
      const r = 60 * Math.cbrt(Math.random())
      pos[i*3] = r * Math.sin(phi) * Math.cos(theta)
      pos[i*3+1] = r * Math.sin(phi) * Math.sin(theta)
      pos[i*3+2] = r * Math.cos(phi)
    }
    return pos
  }, [])

  const color = theme === 'dark' ? '#1e2d45' : '#aaa'

  return (
    <points>
      <bufferGeometry>
        <bufferAttribute
          attach="attributes-position"
          count={800}
          array={positions}
          itemSize={3}
        />
      </bufferGeometry>
      <pointsMaterial size={0.06} color={color} transparent opacity={0.8} sizeAttenuation />
    </points>
  )
}

function GraphScene({ nodes, links, activeId, onSelect, theme }: GraphSceneProps) {
  const [hovered, setHovered] = useState<string | null>(null)
  const controlsRef = useRef<any>(null)
  const isAnyHovered = hovered !== null

  useEffect(() => {
    if (controlsRef.current) {
      controlsRef.current.autoRotate = !isAnyHovered
    }
  }, [isAnyHovered])

  const bgColor = theme === 'dark' ? '#0b0f1a' : '#f5f5f5'

  return (
    <>
      <ambientLight intensity={0.4} />
      <pointLight position={[10, 10, 10]} intensity={0.8} />
      <pointLight position={[-10, -10, -10]} intensity={0.4} color={theme === 'dark' ? '#c8a96e' : '#a0814a'} />
      
      <AllEdges links={links} theme={theme} />
      <EdgeLines links={links} nodes={nodes} activeId={activeId} theme={theme} />
      
      {nodes.map(node => (
        <NodeMesh
          key={node.id}
          node={node}
          isActive={node.id === activeId}
          isHovered={node.id === hovered}
          onHover={setHovered}
          onSelect={() => onSelect(node.id)}
          theme={theme}
        />
      ))}
      
      <Particles theme={theme} />
      
      <OrbitControls
        ref={controlsRef}
        enablePan={false}
        enableZoom={true}
        minDistance={12}
        maxDistance={45}
        autoRotate={true}
        autoRotateSpeed={0.4}
        dampingFactor={0.08}
        enableDamping={true}
      />
    </>
  )
}

function InfoPanel({ activeId, theme }: { activeId: string, theme: 'dark' | 'light' }) {
  const node = nodeMap.get(activeId)
  if (!node) return null

  const [isMobile, setIsMobile] = useState(false)

  useEffect(() => {
    const checkMobile = () => setIsMobile(window.innerWidth < 640)
    checkMobile()
    window.addEventListener('resize', checkMobile)
    return () => window.removeEventListener('resize', checkMobile)
  }, [])

  const connectedLinks = LINKS.filter(l => l.source === activeId || l.target === activeId)
  const connectedIds = connectedLinks.map(l => l.source === activeId ? l.target : l.source)
  const fieldColors = theme === 'dark' ? FIELD_COLORS : FIELD_COLORS_LIGHT
  const fieldColor = fieldColors[node.field]
  const dimFieldColor = fieldColor + '30'

  const styles = {
    panel: {
      background: theme === 'dark' ? '#111827' : '#fff',
      borderTop: `0.5px solid ${theme === 'dark' ? '#1e2d45' : '#ddd'}`,
      padding: isMobile ? '0.5rem' : '0.5rem 1rem'
    },
    label: {
      color: theme === 'dark' ? '#5a6577' : '#888',
      fontFamily: 'monospace', 
      fontSize: '10px', 
      textTransform: 'uppercase'
    },
    name: {
      color: theme === 'dark' ? '#c8a96e' : '#a0814a', 
      fontFamily: 'serif', 
      fontSize: isMobile ? '16px' : '18px', 
      marginTop: '4px'
    },
    badge: {
      background: dimFieldColor, 
      color: fieldColor, 
      fontFamily: 'monospace', 
      fontSize: isMobile ? '10px' : '11px',
      padding: '2px 8px',
      borderRadius: '4px'
    },
    level: {
      color: theme === 'dark' ? '#9fa8b4' : '#666', 
      fontFamily: 'monospace', 
      fontSize: '11px', 
      textTransform: 'uppercase'
    },
    desc: {
      color: theme === 'dark' ? '#e8e2d5' : '#333', 
      fontSize: '12px', 
      lineHeight: '1.5', 
      marginTop: '8px'
    },
    pill: {
      background: theme === 'dark' ? '#c8a96e20' : '#a0814a20', 
      color: theme === 'dark' ? '#c8a96e' : '#a0814a', 
      fontFamily: 'monospace', 
      fontSize: '10px',
      padding: '2px 8px',
      borderRadius: '12px'
    }
  }

  return (
    <div style={styles.panel}>
      <div style={styles.label}>selected concept</div>
      <div style={styles.name}>{node.name}</div>
      <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginTop: '8px', flexWrap: 'wrap' }}>
        <span style={styles.badge}>{node.field}</span>
        <span style={styles.level}>{node.level}</span>
      </div>
      <div style={styles.desc}>{node.description}</div>
      {connectedIds.length > 0 && (
        <div style={{ marginTop: '4px', display: 'flex', gap: '6px', flexWrap: 'wrap' }}>
          {connectedIds.map(id => (
            <span key={id} style={styles.pill}>{id}</span>
          ))}
        </div>
      )}
    </div>
  )
}

function LegendRow({ theme }: { theme: 'dark' | 'light' }) {
  const [isMobile, setIsMobile] = useState(false)

  useEffect(() => {
    const checkMobile = () => setIsMobile(window.innerWidth < 640)
    checkMobile()
    window.addEventListener('resize', checkMobile)
    return () => window.removeEventListener('resize', checkMobile)
  }, [])

  if (isMobile) return null

  const fieldColors = theme === 'dark' ? FIELD_COLORS : FIELD_COLORS_LIGHT
  const linkColors = theme === 'dark' ? LINK_COLORS : LINK_COLORS_LIGHT

  const labelColor = theme === 'dark' ? '#5a6577' : '#888'

  return (
    <div style={{ 
      display: 'flex', 
      gap: '0px', 
      justifyContent: 'center', 
      marginTop: '0px',
      flexWrap: 'wrap',
      padding: '0px'
    }}>
    </div>
  )
}

export default function MathConceptGraph3D({ theme = 'dark' }: MathConceptGraph3DProps) {
  const [activeId, setActiveId] = useState('counting')
  const [isMobile, setIsMobile] = useState(false)

  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth < 640)
    }
    checkMobile()
    window.addEventListener('resize', checkMobile)
    return () => window.removeEventListener('resize', checkMobile)
  }, [])

  const styles = {
    label: {
      color: theme === 'dark' ? '#5a6577' : '#888',
      fontFamily: 'monospace', 
      fontSize: isMobile ? '9px' : '10px', 
      textTransform: 'uppercase', 
      letterSpacing: '0.18em'
    },
    sublabel: {
      color: theme === 'dark' ? '#c8a96e' : '#a0814a',
      fontFamily: 'monospace', 
      fontSize: isMobile ? '10px' : '11px', 
      marginTop: '4px'
    },
    container: {
      height: isMobile ? '320px' : '520px', 
      width: '100%', 
      borderRadius: '8px', 
      overflow: 'hidden'
    }
  }

  return (
    <section className="w-full">
      <div className="text-center mb-2">
        <h3 style={styles.label}>
          live concept graph
        </h3>
      </div>
      <div style={styles.container}>
        <Canvas
          camera={{ position: [0, 0, 28], fov: 60 }}
          style={{ background: 'transparent' }}
          dpr={[1, 2]}
        >
          <color attach="background" args={[theme === 'dark' ? '#0b0f1a' : '#f5f5f5']} />
          <GraphScene
            nodes={NODES}
            links={LINKS}
            activeId={activeId}
            onSelect={setActiveId}
            theme={theme}
          />
        </Canvas>
      </div>
      <InfoPanel activeId={activeId} theme={theme}/>
      {/* <LegendRow theme={theme} /> */}
    </section>
  )
}
